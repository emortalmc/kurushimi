package dev.emortal.api.service.matchmaker;

import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

public record QueueOptions(boolean privateGame, @Nullable String mapId, boolean autoTeleport) {
    public static final QueueOptions DEFAULT = new QueueOptions(false, null, false);

    public static @NotNull Builder builder() {
        return new Builder();
    }

    public static final class Builder {

        private boolean privateGame;
        private @Nullable String mapId;
        private boolean autoTeleport;

        private Builder() {
        }

        public @NotNull Builder privateGame(boolean privateGame) {
            this.privateGame = privateGame;
            return this;
        }

        public @NotNull Builder mapId(@Nullable String mapId) {
            this.mapId = mapId;
            return this;
        }

        public @NotNull Builder autoTeleport(boolean autoTeleport) {
            this.autoTeleport = autoTeleport;
            return this;
        }

        public @NotNull QueueOptions build() {
            return new QueueOptions(this.privateGame, this.mapId, this.autoTeleport);
        }
    }
}
