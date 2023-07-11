package dev.emortal.api.service.matchmaker;

import dev.emortal.api.kurushimi.PendingMatch;
import dev.emortal.api.kurushimi.QueuedPlayer;
import dev.emortal.api.kurushimi.Ticket;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

public record PlayerQueueInfo(@NotNull Ticket ticket, @NotNull QueuedPlayer queuedPlayer, @Nullable PendingMatch pendingMatch) {
}
