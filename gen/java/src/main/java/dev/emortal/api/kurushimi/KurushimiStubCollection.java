package dev.emortal.api.kurushimi;

import dev.emortal.api.service.matchmaker.DefaultMatchmakerService;
import dev.emortal.api.service.matchmaker.MatchmakerService;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;

@SuppressWarnings("OptionalUsedAsFieldOrParameterType")
public class KurushimiStubCollection {
    private static final Logger LOGGER = LoggerFactory.getLogger(KurushimiStubCollection.class);
    private static final boolean DEVELOPMENT = System.getenv("KUBERNETES_SERVICE_HOST") == null;

    private static final @NotNull Optional<MatchmakerService> service;

    static {
        service = createChannel().map(DefaultMatchmakerService::new);
    }

    public static @NotNull Optional<MatchmakerService> getService() {
        return service;
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
