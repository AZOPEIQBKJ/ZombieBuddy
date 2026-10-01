package me.zed_0xff.zombie_buddy.i18n;

import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.Locale;
import java.util.Properties;

/** Approval UI text available before PZ initializes and in standalone child JVMs. */
public final class Messages {
    public static final String LANGUAGE_PROPERTY = "zombiebuddy.language";
    private static final Properties EN = load("en");
    private static final Properties FR = load("fr");

    private Messages() {}

    public static String language() {
        String selected = System.getProperty(LANGUAGE_PROPERTY);
        if (selected == null || selected.isBlank()) selected = Locale.getDefault().getLanguage();
        selected = selected.trim().toLowerCase(Locale.ROOT).replace('_', '-');
        return selected.equals("fr") || selected.startsWith("fr-") ? "fr" : "en";
    }

    public static String childJvmArgument() {
        return "-D" + LANGUAGE_PROPERTY + "=" + language();
    }

    public static String text(String key, Object... args) {
        Properties selected = language().equals("fr") ? FR : EN;
        String template = selected.getProperty(key, EN.getProperty(key, "[" + key + "]"));
        return args.length == 0 ? template : String.format(Locale.ROOT, template, args);
    }

    private static Properties load(String language) {
        Properties properties = new Properties();
        String path = "/me/zed_0xff/zombie_buddy/i18n/messages_" + language + ".properties";
        try (var stream = Messages.class.getResourceAsStream(path)) {
            if (stream == null) throw new IllegalStateException("Missing approval catalogue: " + path);
            properties.load(new InputStreamReader(stream, StandardCharsets.UTF_8));
        } catch (IOException error) {
            throw new IllegalStateException("Unreadable approval catalogue: " + path, error);
        }
        return properties;
    }
}
