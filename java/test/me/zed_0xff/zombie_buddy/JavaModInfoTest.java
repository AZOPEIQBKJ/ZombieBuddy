package me.zed_0xff.zombie_buddy;

import java.io.File;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.AfterEach;
import org.mockito.MockedStatic;
import static org.mockito.Mockito.*;
import static org.junit.jupiter.api.Assertions.*;

import static me.zed_0xff.zombie_buddy.SteamWorkshop.WorkshopItemID;
import java.nio.file.Path;
import java.nio.file.Files;
import org.junit.jupiter.api.io.TempDir;

class JavaModInfoTest {
    @TempDir Path fixture;
    private String HOME_DIR;
    private String CACHE_DIR;

    private MockedStatic<Utils> utilsMock;
    @BeforeEach
    void setUp() throws Exception {
        HOME_DIR = fixture.toString();
        CACHE_DIR = fixture.resolve("Zomboid").toString();
        Path workshop = Path.of(CACHE_DIR, "Workshop/ZBExhume41");
        Files.createDirectories(workshop);
        Files.writeString(workshop.resolve("workshop.txt"), "id=3718604798\n");
        utilsMock = mockStatic(Utils.class, CALLS_REAL_METHODS);
        utilsMock.when(Utils::getCacheDir).thenReturn(CACHE_DIR);
        utilsMock.when(Utils::getCachePath).thenReturn(Path.of(CACHE_DIR));
    }

    @AfterEach
    void tearDown() {
        utilsMock.close();
    }

    @Test
    void workshopItemIdFromInfPath_valid() {
        assertEquals(
                new WorkshopItemID(3718604798L),
                JavaModInfo.workshopItemIdFromInfPath( Path.of(CACHE_DIR, "Workshop/ZBExhume41/Contents/mods/ZBExhume41/common/mod.info"))
                );
        assertEquals(
                new WorkshopItemID(2986022978L),
                JavaModInfo.workshopItemIdFromInfPath( Path.of(HOME_DIR, "/Library/Application Support/Steam/steamapps/workshop/content/108600/2986022978/mods/DoubleDeckerBusInterior"))
                );
    }

    @Test
    void workshopItemIdFromInfPath_invalid() {
        assertNull( JavaModInfo.workshopItemIdFromInfPath( Path.of(CACHE_DIR, "Workshop/ZBExhume41/Contents/mods/ZBExhume41/common/foo/mod.info")));
        assertNull( JavaModInfo.workshopItemIdFromInfPath( Path.of(CACHE_DIR, "Workshop/ZBExhume41/Contents/mods/ZBExhume41/mod.info")));
        assertNull( JavaModInfo.workshopItemIdFromInfPath( Path.of(CACHE_DIR, "Workshop/ZBExhume41/Contents/mods/ZBExhume41/42.13/media/java/ZBExhume41.jar")));
        assertNull( JavaModInfo.workshopItemIdFromInfPath( Path.of("/etc/passwd")));
    }

    @Test
    void rootAdjacentCacheDoesNotHaveAParentFilename() {
        Path root = fixture.toAbsolutePath().getRoot();
        assertNull(JavaModInfo.workshopItemIdFromInfPath(
            root.resolve("Games/Cache/mods/Example/42.21/mod.info")));
        assertNull(JavaModInfo.workshopItemIdFromInfPath(root));
        assertNull(JavaModInfo.workshopItemIdFromInfPath(null));
    }
}
