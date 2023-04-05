package dev.emortal.api.kurushimi;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;

@SuppressWarnings("OptionalUsedAsFieldOrParameterType")
public class KurushimiStubCollection {
    private static final Logger LOGGER = LoggerFactory.getLogger(KurushimiStubCollection.class);
    private static final boolean DEVELOPMENT = System.getenv("HOSTNAME") == null;

    private static final @NotNull Optional<MatchmakerGrpc.MatchmakerFutureStub> futureStub;
    private static final @NotNull Optional<MatchmakerGrpc.MatchmakerStub> stub;
    private static final @NotNull Optional<MatchmakerGrpc.MatchmakerBlockingStub> blockingStub;

    static {
        Optional<ManagedChannel> channel = createChannel();
        if (channel.isPresent()) {
            futureStub = Optional.of(MatchmakerGrpc.newFutureStub(channel.get()));
            stub = Optional.of(MatchmakerGrpc.newStub(channel.get()));
            blockingStub = Optional.of(MatchmakerGrpc.newBlockingStub(channel.get()));
        } else {
            futureStub = Optional.empty();
            stub = Optional.empty();
            blockingStub = Optional.empty();
        }
    }

    public static @NotNull Optional<MatchmakerGrpc.MatchmakerFutureStub> getFutureStub() {
        return futureStub;
    }

    public static @NotNull Optional<MatchmakerGrpc.MatchmakerStub> getStub() {
        return stub;
    }

    public static @NotNull Optional<MatchmakerGrpc.MatchmakerBlockingStub> getBlockingStub() {
        return blockingStub;
    }

    /**
     * @return Optional of a ManagedChannel, empty if service is not enabled.
     */
    private static Optional<ManagedChannel> createChannel() {
        String name = "matchmaker";

        String envPort = System.getenv("MATCHMAKER_SERVICE_PORT");

        if (!DEVELOPMENT) {
            return Optional.of(ManagedChannelBuilder.forAddress(name, Integer.parseInt(envPort))
                    .defaultLoadBalancingPolicy("round_robin")
                    .usePlaintext()
                    .build());
        }

        if (envPort == null) {
            LOGGER.warn("Service {} is not enabled, skipping", name);
            return Optional.empty();
        } else {
            int port = Integer.parseInt(envPort);
            return Optional.of(ManagedChannelBuilder.forAddress("localhost", port)
                    .usePlaintext()
                    .build());
        }
    }
}
