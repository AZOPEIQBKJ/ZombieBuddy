package me.zed_0xff.zombie_buddy;

import static org.junit.jupiter.api.Assertions.*;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;
import me.zed_0xff.zombie_buddy.SteamWorkshop.SteamID64;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

class SignaturePresentationTest {
    @TempDir Path temp;

    @Test void missingSignatureRetainsTypeAndKeepsEnglishNotice() {
        Path jar = temp.resolve("fixture.jar");
        var result = ZBSVerifier.verify(jar, temp.resolve("absent.zbs"), "0".repeat(64));
        assertInstanceOf(ZBSVerifier.MissingSignature.class, result);
        assertTrue(result.detailedMessage.startsWith("Missing .zbs file next to JAR"));
        assertTrue(ZBSVerifier.noticeForUi(result).startsWith("Missing signature file."));
        assertTrue(ZBSVerifier.noticeForUi(result).contains("next to JAR"));
    }

    @Test void uploaderMismatchRemainsInvalidWithEnglishDetails() throws Exception {
        Path jar = temp.resolve("fixture.jar"), sidecar = temp.resolve("fixture.jar.zbs");
        Files.writeString(sidecar, "ZBS\nSteamID64:76561198000000000\nSignature:" + "0".repeat(128) + "\n");
        var result = ZBSVerifier.verify(jar, sidecar, "0".repeat(64), new SteamID64(76561198000000001L), Map.of());
        assertInstanceOf(ZBSVerifier.InvalidSignature.class, result);
        assertEquals("Declared SteamID64 does not match Workshop item uploader.", result.detailedMessage);
        assertTrue(ZBSVerifier.noticeForUi(result).contains("does not match"));
    }

    @Test void missingNotAllowedKeepsBlockingPolicyAndEnglishNotice() {
        var result = ZBSVerifier.check(temp.resolve("fixture.jar"), "0".repeat(64), null, false, false, Map.of());
        assertFalse(result.flags().has(ModFlags.MF_VALID));
        assertEquals("missing .zbs file; allow_unsigned_mods=false", result.blockReason());
        assertEquals("Missing .zbs file (allow_unsigned_mods=false)", result.notice());
    }
}
