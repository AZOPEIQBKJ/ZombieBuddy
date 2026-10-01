package me.zed_0xff.zombie_buddy.i18n;

import static org.junit.jupiter.api.Assertions.*;
import java.util.Locale;
import org.junit.jupiter.api.Test;

class MessagesTest {
    @Test void interfaceStaysEnglishWithNonEnglishJvmLocale() {
        Locale previous = Locale.getDefault();
        String legacyOverride = System.getProperty("zombiebuddy.language");
        try {
            Locale.setDefault(Locale.FRANCE);
            System.setProperty("zombiebuddy.language", "fr");
            assertEquals("Yes", Messages.text("answer.yes"));
            assertEquals("Cancel", Messages.text("answer.cancel"));
        } finally {
            Locale.setDefault(previous);
            if (legacyOverride == null) System.clearProperty("zombiebuddy.language");
            else System.setProperty("zombiebuddy.language", legacyOverride);
        }
    }

    @Test void formattingPreservesTechnicalDetailsAndLineBreaks() {
        String detail = "D:\\mods\\a%20.jar";
        assertTrue(Messages.text("signature.unreadable", detail).contains(detail));
        assertTrue(Messages.text("tiny.invalid", "detail").contains("\n\ndetail\n\n"));
    }

    @Test void unknownKeyIsVisibleAndFrenchCatalogueIsNotDistributed() {
        assertEquals("[not.a.key]", Messages.text("not.a.key"));
        assertNull(Messages.class.getResource("messages_fr.properties"));
    }
}
