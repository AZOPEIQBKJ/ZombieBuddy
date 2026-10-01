package me.zed_0xff.zombie_buddy.frontend;

import static org.junit.jupiter.api.Assertions.*;
import static me.zed_0xff.zombie_buddy.ModFlags.MF_PERSIST;
import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.io.StringReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.util.List;
import me.zed_0xff.zombie_buddy.JarBatchApprovalProtocol.Entry;
import me.zed_0xff.zombie_buddy.ModFlags;
import org.junit.jupiter.api.Test;

class ConsoleApprovalTest {

    @Test void yesPersistsOnlyWhenExplicitlyChosen() {
        var entry = entry(Entry.ZBSignature.none());
        String output = run("yes\nyes\n", entry);
        assertTrue(entry.decision);
        assertTrue(entry.flags.has(MF_PERSIST));
        assertTrue(output.contains("Allow this Java mod"));
        assertTrue(output.contains("[y/n]"));
    }

    @Test void noAndEofNeverGrantApproval() {
        var denied = entry(Entry.ZBSignature.none());
        run("no\nno\n", denied);
        assertFalse(denied.decision);
        var eof = entry(Entry.ZBSignature.none());
        run("", eof);
        assertFalse(eof.decision);
        assertFalse(eof.flags.has(MF_PERSIST));
    }

    @Test void invalidSignatureCannotBeApprovedWithYes() {
        var invalid = entry(new Entry.ZBSignature(false, null, "Invalid signature."));
        String output = run("yes\n", invalid);
        assertFalse(invalid.decision);
        assertFalse(output.contains("Allow this Java mod"));
        assertTrue(output.contains("loading is denied"));
    }

    @Test void englishYesAndRetryKeepExistingAnswerContract() {
        var entry = entry(Entry.ZBSignature.none());
        String output = run("o\nyes\nno\n", entry);
        assertTrue(entry.decision);
        assertFalse(entry.flags.has(MF_PERSIST));
        assertTrue(output.contains("Please answer y or n."));
    }

    private String run(String input, Entry entry) {
        var bytes = new ByteArrayOutputStream();
        var frontend = new ConsoleModApprovalFrontend(new StringReader(input), new PrintStream(bytes, true, StandardCharsets.UTF_8));
        frontend.approvePendingMods(List.of(entry));
        return bytes.toString(StandardCharsets.UTF_8);
    }

    private Entry entry(Entry.ZBSignature signature) {
        return new Entry("fixture", null, Path.of("fixture.jar"), Path.of("mod.info"),
            "a".repeat(64), null, null, ModFlags.EMPTY, "Fixture", signature, null);
    }
}
