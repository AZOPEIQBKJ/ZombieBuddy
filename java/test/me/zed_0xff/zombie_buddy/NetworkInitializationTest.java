package me.zed_0xff.zombie_buddy;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import java.io.IOException;
import java.lang.reflect.InvocationTargetException;
import java.net.http.HttpClient;
import java.nio.file.Path;
import java.util.Set;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class NetworkInitializationTest {
    @TempDir Path config;

    @Test
    void failedClientCreationUsesExistingUnknownAndUnverifiablePaths() throws Exception {
        String previous = Agent.arguments.put("config_dir", config.toString());
        try (var http = mockStatic(HttpClient.class)) {
            http.when(HttpClient::newBuilder).thenThrow(new IllegalStateException("test: selector unavailable"));
            var id = new SteamWorkshop.WorkshopItemID(1234567890L);
            var details = SteamWorkshop.fetchItemDetails(Set.of(id));
            assertNull(details.get(id).ban().status(), "network failure must not become a known clean status");
            assertNull(details.get(id).creatorSteamId64());
            assertTrue(KnownAuthors.loadAuthors().isEmpty(), "no unsigned substitute author registry");

            var method = ZBSVerifier.class.getDeclaredMethod("fetchJavaModZBSHexesFromSteam", SteamWorkshop.SteamID64.class);
            method.setAccessible(true);
            var failure = assertThrows(InvocationTargetException.class,
                () -> method.invoke(null, new SteamWorkshop.SteamID64(76561198000000000L)));
            assertInstanceOf(IOException.class, failure.getCause());
        } finally {
            if (previous == null) Agent.arguments.remove("config_dir");
            else Agent.arguments.put("config_dir", previous);
        }
    }
}
