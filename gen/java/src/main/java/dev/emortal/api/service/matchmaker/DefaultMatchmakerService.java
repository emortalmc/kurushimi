package dev.emortal.api.service.matchmaker;

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Message;
import com.google.rpc.Status;
import dev.emortal.api.kurushimi.ChangePlayerMapVoteErrorResponse;
import dev.emortal.api.kurushimi.ChangePlayerMapVoteRequest;
import dev.emortal.api.kurushimi.DequeueByPlayerErrorResponse;
import dev.emortal.api.kurushimi.DequeueByPlayerRequest;
import dev.emortal.api.kurushimi.GetPlayerQueueInfoRequest;
import dev.emortal.api.kurushimi.GetPlayerQueueInfoResponse;
import dev.emortal.api.kurushimi.MatchmakerGrpc;
import dev.emortal.api.kurushimi.QueueByPlayerErrorResponse;
import dev.emortal.api.kurushimi.QueueByPlayerRequest;
import dev.emortal.api.kurushimi.QueueInitialLobbyByPlayerRequest;
import dev.emortal.api.kurushimi.SendPlayerToLobbyRequest;
import io.grpc.ManagedChannel;
import io.grpc.StatusRuntimeException;
import io.grpc.protobuf.StatusProto;
import org.jetbrains.annotations.ApiStatus;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collection;
import java.util.UUID;

@ApiStatus.Internal
public final class DefaultMatchmakerService implements MatchmakerService {
    private static final Logger LOGGER = LoggerFactory.getLogger(DefaultMatchmakerService.class);

    private final MatchmakerGrpc.MatchmakerBlockingStub grpc;

    @ApiStatus.Internal
    public DefaultMatchmakerService(@NotNull ManagedChannel channel) {
        this.grpc = MatchmakerGrpc.newBlockingStub(channel);
    }

    @Override
    public @Nullable PlayerQueueInfo getPlayerQueueInfo(@NotNull UUID playerId) {
        var request = GetPlayerQueueInfoRequest.newBuilder().setPlayerId(playerId.toString()).build();

        GetPlayerQueueInfoResponse response;
        try {
            response = this.grpc.getPlayerQueueInfo(request);
        } catch (StatusRuntimeException exception) {
            if (exception.getStatus() == io.grpc.Status.NOT_FOUND) return null;
            throw exception;
        }

        return new PlayerQueueInfo(response.getTicket(), response.getQueuedPlayer(), response.hasPendingMatch() ? response.getPendingMatch() : null);
    }

    @Override
    public @NotNull QueuePlayerResult queuePlayer(@NotNull String gameModeId, @NotNull UUID playerId, @NotNull QueueOptions options) {
        var requestBuilder = QueueByPlayerRequest.newBuilder().setGameModeId(gameModeId).setPlayerId(playerId.toString());
        if (options.privateGame()) requestBuilder.setPrivateGame(true);
        if (options.mapId() != null) requestBuilder.setMapId(options.mapId());
        if (options.autoTeleport()) requestBuilder.setAutoTeleport(true);
        var request = requestBuilder.build();

        try {
            this.grpc.queueByPlayer(request);
            return QueuePlayerResult.SUCCESS;
        } catch (StatusRuntimeException exception) {
            var error = getError(exception, QueueByPlayerErrorResponse.class);

            return switch (error.getReason()) {
                case ALREADY_IN_QUEUE -> QueuePlayerResult.ALREADY_IN_QUEUE;
                case INVALID_GAME_MODE -> QueuePlayerResult.INVALID_GAME_MODE;
                case GAME_MODE_DISABLED -> QueuePlayerResult.GAME_MODE_DISABLED;
                case INVALID_MAP -> QueuePlayerResult.INVALID_MAP;
                case PARTY_TOO_LARGE -> QueuePlayerResult.PARTY_TOO_LARGE;
                case NO_PERMISSION -> QueuePlayerResult.NO_PERMISSION;
                case UNRECOGNIZED -> throw exception;
            };
        }
    }

    @Override
    public @NotNull DequeuePlayerResult dequeuePlayer(@NotNull UUID playerId) {
        var request = DequeueByPlayerRequest.newBuilder().setPlayerId(playerId.toString()).build();

        try {
            this.grpc.dequeueByPlayer(request);
            return DequeuePlayerResult.SUCCESS;
        } catch (StatusRuntimeException exception) {
            var error = getError(exception, DequeueByPlayerErrorResponse.class);

            return switch (error.getReason()) {
                case NOT_IN_QUEUE -> DequeuePlayerResult.NOT_IN_QUEUE;
                case NO_PERMISSION -> DequeuePlayerResult.NO_PERMISSION;
                case ALREADY_MARKED_FOR_DEQUEUE -> DequeuePlayerResult.ALREADY_MARKED_FOR_DEQUEUE;
                case UNRECOGNIZED -> throw exception;
            };
        }
    }

    @Override
    public void sendPlayerToLobby(@NotNull UUID playerId, boolean sendParty) {
        var requestBuilder = SendPlayerToLobbyRequest.newBuilder().addPlayerIds(playerId.toString());
        if (!sendParty) requestBuilder.setSendParties(false);
        var request = requestBuilder.build();

        this.grpc.sendPlayersToLobby(request);
    }

    @Override
    public void sendPlayersToLobby(@NotNull Collection<UUID> playerIds, boolean sendParties) {
        var requestBuilder = SendPlayerToLobbyRequest.newBuilder();
        for (UUID id : playerIds) {
            requestBuilder.addPlayerIds(id.toString());
        }
        if (!sendParties) requestBuilder.setSendParties(false);
        var request = requestBuilder.build();

        this.grpc.sendPlayersToLobby(request);
    }

    @Override
    public void queueInitialLobby(@NotNull UUID playerId) {
        var request = QueueInitialLobbyByPlayerRequest.newBuilder().setPlayerId(playerId.toString()).build();

        this.grpc.queueInitialLobbyByPlayer(request);
    }

    @Override
    public @NotNull ChangeMapVoteResult changeMapVote(@NotNull UUID playerId, @NotNull String newMapId) {
        var request = ChangePlayerMapVoteRequest.newBuilder()
                .setPlayerId(playerId.toString())
                .setMapId(newMapId)
                .build();

        try {
            this.grpc.changePlayerMapVote(request);
            return ChangeMapVoteResult.SUCCESS;
        } catch (StatusRuntimeException exception) {
            var error = getError(exception, ChangePlayerMapVoteErrorResponse.class);

            return switch (error.getReason()) {
                case NOT_IN_QUEUE -> ChangeMapVoteResult.NOT_IN_QUEUE;
                case INVALID_MAP -> ChangeMapVoteResult.INVALID_MAP;
                case UNRECOGNIZED -> throw exception;
            };
        }
    }

    private static <T extends Message> @NotNull T getError(@NotNull StatusRuntimeException exception, @NotNull Class<T> errorType) {
        Status status = StatusProto.fromStatusAndTrailers(exception.getStatus(), exception.getTrailers());
        if (status.getDetailsCount() == 0) {
            LOGGER.error("No error details provided in response to AddRoleToPlayer from server!", exception);
            throw new IllegalStateException("No error details provided in response to AddRoleToPlayer from server!");
        }

        try {
            return status.getDetails(0).unpack(errorType);
        } catch (InvalidProtocolBufferException invalidException) {
            LOGGER.error("Error details provided in response to AddRoleToPlayer from server were invalid!", exception);
            throw new IllegalStateException("Error details provided in response to AddRoleToPlayer from server were invalid!", invalidException);
        }
    }
}
