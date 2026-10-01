package me.zed_0xff.zombie_buddy.i18n;

import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.Locale;
import java.util.Properties;

/** English approval text available before PZ initializes and in standalone child JVMs. */
public final class Messages {
    private static final Properties EN = load("en");

    private Messages() {}

    public static String text(String key, Object... args) {
        String template = EN.getProperty(key, "[" + key + "]");
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
