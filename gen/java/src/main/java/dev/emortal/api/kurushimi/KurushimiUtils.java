package dev.emortal.api.kurushimi;

import dev.emortal.api.kurushimi.messages.MatchCreatedMessage;
import dev.emortal.api.kurushimi.messages.PendingMatchCreatedMessage;
import dev.emortal.api.kurushimi.messages.PendingMatchDeletedMessage;
import dev.emortal.api.kurushimi.messages.PendingMatchUpdatedMessage;
import dev.emortal.api.kurushimi.messages.TicketCreatedMessage;
import dev.emortal.api.kurushimi.messages.TicketDeletedMessage;
import dev.emortal.api.kurushimi.messages.TicketUpdatedMessage;
import dev.emortal.api.utils.parser.ProtoParserRegistry;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@SuppressWarnings("unused")
public class KurushimiUtils {
    private static final Logger LOGGER = LoggerFactory.getLogger(KurushimiUtils.class);

    public static void registerParserRegistry() {
        LOGGER.info("Registering Kurushimi parser registry");

        ProtoParserRegistry.register(TicketCreatedMessage.getDefaultInstance(), TicketCreatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.register(TicketUpdatedMessage.getDefaultInstance(), TicketUpdatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.register(TicketDeletedMessage.getDefaultInstance(), TicketDeletedMessage::parseFrom, "matchmaker");

        ProtoParserRegistry.register(PendingMatchCreatedMessage.getDefaultInstance(), PendingMatchCreatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.register(PendingMatchUpdatedMessage.getDefaultInstance(), PendingMatchUpdatedMessage::parseFrom, "matchmaker");
        ProtoParserRegistry.register(PendingMatchDeletedMessage.getDefaultInstance(), PendingMatchDeletedMessage::parseFrom, "matchmaker");

        ProtoParserRegistry.register(MatchCreatedMessage.getDefaultInstance(), MatchCreatedMessage::parseFrom, "matchmaker");
    }
}
