package me.zed_0xff.zombie_buddy;

import java.io.IOException;
import java.net.http.HttpClient;
import java.util.Objects;
import java.util.function.Supplier;

/** Keep selector/loopback creation failures inside the existing network error paths. */
final class LazyHttpClient {
    private final Supplier<HttpClient> factory;
    private volatile HttpClient client;

    LazyHttpClient(Supplier<HttpClient> factory) {
        this.factory = Objects.requireNonNull(factory);
    }

    HttpClient get() throws IOException {
        HttpClient current = client;
        if (current != null) return current;
        synchronized (this) {
            if (client == null) {
                try {
                    client = Objects.requireNonNull(factory.get());
                } catch (RuntimeException e) {
                    // Leave client unset so a later request may retry after a transient failure.
                    throw new IOException("HTTP client unavailable", e);
                }
            }
            return client;
        }
    }
}
