package me.zed_0xff.zombie_buddy.i18n;

import static org.junit.jupiter.api.Assertions.*;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.Locale;
import java.util.Properties;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

class MessagesTest {
    private final String previous = System.getProperty(Messages.LANGUAGE_PROPERTY);
    private final Locale previousLocale = Locale.getDefault();

    @AfterEach void restore() {
        if (previous == null) System.clearProperty(Messages.LANGUAGE_PROPERTY);
        else System.setProperty(Messages.LANGUAGE_PROPERTY, previous);
        Locale.setDefault(previousLocale);
    }

    @Test void frenchRegionalLocaleAndExplicitOverride() {
        System.clearProperty(Messages.LANGUAGE_PROPERTY);
        Locale.setDefault(Locale.forLanguageTag("fr-BE"));
        assertEquals("Oui", Messages.text("answer.yes"));
        System.setProperty(Messages.LANGUAGE_PROPERTY, "en");
        assertEquals("Yes", Messages.text("answer.yes"));
    }

    @Test void normalizedChildLanguageDoesNotInjectJvmArguments() {
        System.setProperty(Messages.LANGUAGE_PROPERTY, " FR_be ");
        assertEquals("-Dzombiebuddy.language=fr", Messages.childJvmArgument());
        System.setProperty(Messages.LANGUAGE_PROPERTY, "en -javaagent:other.jar");
        assertEquals("-Dzombiebuddy.language=en", Messages.childJvmArgument());
    }

    @Test void unsupportedLanguageFallsBackToEnglish() {
        System.setProperty(Messages.LANGUAGE_PROPERTY, "de");
        assertEquals("Cancel", Messages.text("answer.cancel"));
        assertEquals("[not.a.key]", Messages.text("not.a.key"));
    }

    @Test void unicodeFormattingPreservesApostrophesPathsAndLineBreaks() {
        System.setProperty(Messages.LANGUAGE_PROPERTY, "fr");
        assertTrue(Messages.text("column.updated").contains("à"));
        String notice = Messages.text("signature.unreadable", "D:\\mods\\a%20.jar");
        assertTrue(notice.contains("D:\\mods\\a%20.jar"));
        assertTrue(Messages.text("approval.intro").contains("d'autoriser"));
        assertTrue(Messages.text("tiny.invalid", "detail").contains("\n\ndetail\n\n"));
    }

    @Test void cataloguesHaveMatchingKeysAndFormatArguments() throws Exception {
        Properties en = catalogue("en"), fr = catalogue("fr");
        assertEquals(en.keySet(), fr.keySet());
        assertTrue(en.size() >= 60);
        for (String key : en.stringPropertyNames()) {
            assertFalse(fr.getProperty(key).isBlank(), key);
            long count = en.getProperty(key).chars().filter(c -> c == '%').count();
            assertEquals(count, fr.getProperty(key).chars().filter(c -> c == '%').count(), key);
            Object[] args = new Object[(int) count];
            java.util.Arrays.fill(args, "fixture");
            assertDoesNotThrow(() -> String.format(Locale.ROOT, fr.getProperty(key), args), key);
            // ImGui's existing default font covers Latin-1; avoid adding unsupported French glyphs.
            assertTrue(fr.getProperty(key).chars().allMatch(c -> c <= 255), key);
        }
    }

    private Properties catalogue(String language) throws Exception {
        Properties properties = new Properties();
        try (var stream = Messages.class.getResourceAsStream("messages_" + language + ".properties")) {
            assertNotNull(stream);
            properties.load(new InputStreamReader(stream, StandardCharsets.UTF_8));
        }
        return properties;
    }
}
