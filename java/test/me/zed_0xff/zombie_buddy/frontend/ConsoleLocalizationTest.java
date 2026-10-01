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
import me.zed_0xff.zombie_buddy.i18n.Messages;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

class ConsoleLocalizationTest {
    private final String previous = System.getProperty(Messages.LANGUAGE_PROPERTY);
    @AfterEach void restore() {
        if (previous == null) System.clearProperty(Messages.LANGUAGE_PROPERTY);
        else System.setProperty(Messages.LANGUAGE_PROPERTY, previous);
    }

    @Test void frenchOuiPersistsOnlyWhenExplicitlyChosen() {
        var entry = entry(Entry.ZBSignature.none());
        String output = run("fr", "oui\noui\n", entry);
        assertTrue(entry.decision);
        assertTrue(entry.flags.has(MF_PERSIST));
        assertTrue(output.contains("Autoriser le chargement"));
        assertTrue(output.contains("[o/n]"));
    }

    @Test void frenchNonAndEofNeverGrantApproval() {
        var denied = entry(Entry.ZBSignature.none());
        run("fr", "non\nnon\n", denied);
        assertFalse(denied.decision);
        var eof = entry(Entry.ZBSignature.none());
        run("fr", "", eof);
        assertFalse(eof.decision);
        assertFalse(eof.flags.has(MF_PERSIST));
    }

    @Test void invalidSignatureCannotBeApprovedWithOui() {
        var invalid = entry(new Entry.ZBSignature(false, null, "Signature invalide."));
        String output = run("fr", "oui\n", invalid);
        assertFalse(invalid.decision);
        assertFalse(output.contains("Autoriser le chargement"));
        assertTrue(output.contains("chargement est refusé"));
    }

    @Test void englishYesAndRetryKeepExistingAnswerContract() {
        var entry = entry(Entry.ZBSignature.none());
        String output = run("en", "o\nyes\nno\n", entry);
        assertTrue(entry.decision);
        assertFalse(entry.flags.has(MF_PERSIST));
        assertTrue(output.contains("Please answer y or n."));
    }

    private String run(String language, String input, Entry entry) {
        System.setProperty(Messages.LANGUAGE_PROPERTY, language);
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
