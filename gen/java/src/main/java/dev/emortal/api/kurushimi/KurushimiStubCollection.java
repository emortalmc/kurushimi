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

    private static final @NotNull Optional<FrontendGrpc.FrontendFutureStub> futureStub;
    private static final @NotNull Optional<FrontendGrpc.FrontendStub> stub;
    private static final @NotNull Optional<FrontendGrpc.FrontendBlockingStub> blockingStub;

    static {
        Optional<ManagedChannel> channel = createChannel();
        if (channel.isPresent()) {
            futureStub = Optional.of(FrontendGrpc.newFutureStub(channel.get()));
            stub = Optional.of(FrontendGrpc.newStub(channel.get()));
            blockingStub = Optional.of(FrontendGrpc.newBlockingStub(channel.get()));
        } else {
            futureStub = Optional.empty();
            stub = Optional.empty();
            blockingStub = Optional.empty();
        }
    }

    public static @NotNull Optional<FrontendGrpc.FrontendFutureStub> getFutureStub() {
        return futureStub;
    }

    public static @NotNull Optional<FrontendGrpc.FrontendStub> getStub() {
        return stub;
    }

    public static @NotNull Optional<FrontendGrpc.FrontendBlockingStub> getBlockingStub() {
        return blockingStub;
    }

    /**
     * @return Optional of a ManagedChannel, empty if service is not enabled.
     */
    private static Optional<ManagedChannel> createChannel() {
        String name = "matchmaker";

        if (!DEVELOPMENT) {
            return Optional.of(ManagedChannelBuilder.forAddress(name, 9090)
                    .defaultLoadBalancingPolicy("round_robin")
                    .usePlaintext()
                    .build());
        }
        String envPort = System.getenv(name + "-port"); // only used for dev

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
