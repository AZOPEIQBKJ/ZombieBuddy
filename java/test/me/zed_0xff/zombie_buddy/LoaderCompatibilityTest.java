package me.zed_0xff.zombie_buddy;

import org.junit.jupiter.api.Test;
import java.util.ArrayList;
import java.util.List;
import static org.junit.jupiter.api.Assertions.*;

class LoaderCompatibilityTest {
    @Test
    void legacyBinaryDescriptorRemainsAvailableAlongsideList() throws Exception {
        assertEquals(void.class, Loader.class.getMethod("loadMods", ArrayList.class).getReturnType());
        assertEquals(void.class, Loader.class.getMethod("loadMods", List.class).getReturnType());
    }
}
