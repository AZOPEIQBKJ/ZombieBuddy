package me.zed_0xff.zombie_buddy;

import org.junit.jupiter.api.Test;
import java.io.IOException;
import java.net.http.HttpClient;
import java.util.ArrayList;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.atomic.AtomicInteger;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class LazyHttpClientTest {
    @Test
    void constructionIsDeferredAndTransientFailuresCanRetry() throws Exception {
        AtomicInteger attempts = new AtomicInteger();
        HttpClient expected = mock(HttpClient.class);
        IllegalStateException failure = new IllegalStateException("loopback unavailable");
        LazyHttpClient client = new LazyHttpClient(() -> {
            if (attempts.incrementAndGet() == 1) throw failure;
            return expected;
        });
        assertEquals(0, attempts.get());
        assertSame(failure, assertThrows(IOException.class, client::get).getCause());
        assertSame(expected, client.get());
        assertSame(expected, client.get());
        assertEquals(2, attempts.get());
    }

    @Test
    void concurrentRequestsShareOneSuccessfullyConstructedClient() throws Exception {
        AtomicInteger attempts = new AtomicInteger();
        HttpClient expected = mock(HttpClient.class);
        LazyHttpClient client = new LazyHttpClient(() -> {
            attempts.incrementAndGet();
            return expected;
        });
        try (var executor = Executors.newFixedThreadPool(4)) {
            var futures = new ArrayList<Future<HttpClient>>();
            for (int i = 0; i < 16; i++) futures.add(executor.submit(client::get));
            for (var future : futures) assertSame(expected, future.get());
        }
        assertEquals(1, attempts.get());
    }
}
