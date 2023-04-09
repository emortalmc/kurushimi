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

    public static void registerParserRegistry() {
        ProtoParserRegistry.registerKafka(TicketCreatedMessage.getDefaultInstance(), TicketCreatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(TicketUpdatedMessage.getDefaultInstance(), TicketUpdatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(TicketDeletedMessage.getDefaultInstance(), TicketDeletedMessage::parseFrom, "matchmaker");

        ProtoParserRegistry.registerKafka(PendingMatchCreatedMessage.getDefaultInstance(), PendingMatchCreatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(PendingMatchUpdatedMessage.getDefaultInstance(), PendingMatchUpdatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.registerKafka(PendingMatchDeletedMessage.getDefaultInstance(), PendingMatchDeletedMessage::parseFrom, "matchmaker");

        ProtoParserRegistry.registerKafka(MatchCreatedMessage.getDefaultInstance(), MatchCreatedMessage::parseFrom, "matchmaker");
    }
}
