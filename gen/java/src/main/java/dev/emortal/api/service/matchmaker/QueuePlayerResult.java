package dev.emortal.api.service.matchmaker;

public enum QueuePlayerResult {

    SUCCESS,
    ALREADY_IN_QUEUE,
    INVALID_GAME_MODE,
    GAME_MODE_DISABLED,
    INVALID_MAP,
    PARTY_TOO_LARGE,
    NO_PERMISSION
}
