package dev.emortal.api.service.matchmaker;

public enum DequeuePlayerResult {

    SUCCESS,
    NOT_IN_QUEUE,
    NO_PERMISSION,
    ALREADY_MARKED_FOR_DEQUEUE
}
