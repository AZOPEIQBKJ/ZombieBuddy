package me.zed_0xff.zombie_buddy;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import java.nio.file.Files;
import java.nio.file.Path;
import static org.junit.jupiter.api.Assertions.*;

class CommunityUpdatesTest {
    @TempDir Path folder;

    @Test
    void previewNeverOverwritesOrStagesAnAutomaticReplacement() throws Exception {
        Path current = folder.resolve("ZombieBuddy.jar");
        Path offered = folder.resolve("offered.jar");
        Files.writeString(current, "community");
        Files.writeString(offered, "upstream or untrusted");
        for (String version : new String[]{"2.3.3", "3.0.0", "999.0.0-community.1"}) {
            assertFalse(SelfUpdater.checkAndUpdateIfNewer(offered, current, "2.3.3-community.1", 0));
            SelfUpdater.performUpdate(current, offered, version);
            assertEquals("community", Files.readString(current));
            assertFalse(Files.exists(folder.resolve("ZombieBuddy.jar.new")));
            assertFalse(Files.exists(folder.resolve("ZombieBuddy.jar.bak")));
            assertNull(SelfUpdater.getNewVersion());
        }
    }
}
