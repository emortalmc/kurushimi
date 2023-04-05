package dev.emortal.api.kurushimi;

import com.google.common.util.concurrent.Futures;
import dev.emortal.api.kurushimi.messages.MatchCreatedMessage;
import dev.emortal.api.kurushimi.messages.PendingMatchCreatedMessage;
import dev.emortal.api.kurushimi.messages.PendingMatchDeletedMessage;
import dev.emortal.api.kurushimi.messages.PendingMatchUpdatedMessage;
import dev.emortal.api.kurushimi.messages.TicketCreatedMessage;
import dev.emortal.api.kurushimi.messages.TicketDeletedMessage;
import dev.emortal.api.kurushimi.messages.TicketUpdatedMessage;
import dev.emortal.api.utils.callback.FunctionalFutureCallback;
import dev.emortal.api.utils.parser.ProtoParserRegistry;
import net.minestom.server.MinecraftServer;
import net.minestom.server.entity.Player;
import net.minestom.server.event.EventFilter;
import net.minestom.server.event.EventNode;
import net.minestom.server.event.player.PlayerDisconnectEvent;
import net.minestom.server.event.trait.PlayerEvent;
import net.minestom.server.timer.Task;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.temporal.ChronoUnit;
import java.util.Collection;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;
import java.util.concurrent.ForkJoinPool;
import java.util.concurrent.atomic.AtomicBoolean;

public class KurushimiUtils {
    private static final Logger LOGGER = LoggerFactory.getLogger(KurushimiUtils.class);
    private static final EventNode<PlayerEvent> EVENT_NODE = EventNode.type("kurushimi-utils", EventFilter.PLAYER);

    static {
        MinecraftServer.getGlobalEventHandler().addChild(EVENT_NODE);
    }

    public static void registerParserRegistry() {
        ProtoParserRegistry.registerKafka(TicketCreatedMessage.getDefaultInstance(), TicketCreatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(TicketUpdatedMessage.getDefaultInstance(), TicketUpdatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(TicketDeletedMessage.getDefaultInstance(), TicketDeletedMessage::parseFrom, "matchmaker");

        ProtoParserRegistry.registerKafka(PendingMatchCreatedMessage.getDefaultInstance(), PendingMatchCreatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(PendingMatchUpdatedMessage.getDefaultInstance(), PendingMatchUpdatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(PendingMatchDeletedMessage.getDefaultInstance(), PendingMatchDeletedMessage::parseFrom, "matchmaker");

        ProtoParserRegistry.registerKafka(MatchCreatedMessage.getDefaultInstance(), MatchCreatedMessage::parseFrom, "matchmaker");
    }

    /**
     * Note: The failure runnable is only run at the end of the time if players are not sent.
     * If there are other errors, they may only affect one player and resolve with retries.
     *
     * @param players         The player ids to queue for a lobby
     * @param successRunnable A runnable to run when all players are connected to the lobby.
     * @param failureRunnable A runnable to run when the sender gives up sending players.
     * @param retries         The amount of retries to send players to the lobby before giving up.
     */
    // todo retries
    // todo store player tickets so if we assume a cancellation, we delete the ticket
    public static void sendToLobby(Collection<? extends Player> players, Runnable successRunnable,
                                   Runnable failureRunnable, int retries) {
        if (KurushimiStubCollection.getStub().isEmpty()) {
            throw new IllegalStateException("Kurushimi stub is not present.");
        }

        Set<? extends Player> remainingPlayers = new HashSet<>(players);
        AtomicBoolean finished = new AtomicBoolean(false);

        EventNode<PlayerEvent> localNode = EventNode.type(UUID.randomUUID().toString(), EventFilter.PLAYER,
                (event, player) -> players.contains(player));
        EVENT_NODE.addChild(localNode);

        Task task = MinecraftServer.getSchedulerManager().buildTask(() -> {
            boolean shouldRun = finished.compareAndSet(false, true);
            if (shouldRun) {
                failureRunnable.run();
                EVENT_NODE.removeChild(localNode);
            }
        }).delay(10, ChronoUnit.SECONDS).schedule();

        localNode.addListener(PlayerDisconnectEvent.class, event -> {
            remainingPlayers.remove(event.getPlayer());
            if (remainingPlayers.isEmpty()) {
                task.cancel();
                EVENT_NODE.removeChild(localNode);
                successRunnable.run();
            }
        });

        for (Player player : players) {
            sendToLobby(player, () -> {
                // failure
                LOGGER.warn("Failed to create ticket to send player {} to lobby.", player.getUsername());
            });
        }
    }

    private static void sendToLobby(@NotNull Player player, @NotNull Runnable failureRunnable) {
        var ticketFuture = KurushimiStubCollection.getFutureStub().get().queueByPlayer(QueueByPlayerRequest.newBuilder()
                .setPlayerId(player.getUuid().toString())
                .setGameModeId("lobby")
                .build());

        Futures.addCallback(ticketFuture, FunctionalFutureCallback.create(
                ticket -> {
                }, // Do nothing. We simply detect if the player gets teleported
                throwable -> {
                    failureRunnable.run();
                    // todo log
                }
        ), ForkJoinPool.commonPool());
    }

    /**
     * @param players         The player ids to queue for a lobby
     * @param successRunnable A runnable to run when all players are connected to the lobby.
     * @param failureRunnable A runnable to run when the sender gives up sending players.
     */
    public static void sendToLobby(Collection<? extends Player> players, Runnable successRunnable,
                                   Runnable failureRunnable) {
        sendToLobby(players, successRunnable, failureRunnable, 1);
    }
}
