package me.zed_0xff.zombie_buddy.frontend;

import static me.zed_0xff.zombie_buddy.ModFlags.MF_PERSIST;
import static me.zed_0xff.zombie_buddy.i18n.Messages.text;

import me.zed_0xff.zombie_buddy.*;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.Reader;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;
import java.util.Locale;

/**
 * Text-mode approvals on {@link System#in} / {@link System#out}.
 * Intended for headless dedicated servers where Swing/TinyFD are unavailable or undesirable.
 */
public final class ConsoleModApprovalFrontend implements ModApprovalFrontend {
    private static final String DATE_FORMAT = "yyyy-MM-dd";

    private final BufferedReader in;
    private final PrintStream out;

    public ConsoleModApprovalFrontend() {
        this(new InputStreamReader(System.in, StandardCharsets.UTF_8), System.out);
    }

    ConsoleModApprovalFrontend(Reader input, PrintStream output) {
        in = new BufferedReader(input);
        out = output;
    }

    @Override
    public List<JarBatchApprovalProtocol.Entry> approvePendingMods(List<JarBatchApprovalProtocol.Entry> pending) {
        if (pending.isEmpty()) {
            return pending;
        }
        out.println(text("console.intro", pending.size()));
        for (JarBatchApprovalProtocol.Entry e : pending) {
            out.println();
            out.println("---");
            out.println(text("console.mod", e.modId));
            out.println(text("console.workshop", e.workshopItemId != null ? e.workshopItemId.value() : text("status.none")));
            out.println(text("console.jar", e.jarAbsolutePath));
            out.println(text("console.hash", e.sha256));
            out.println(text("console.updated", formatDate(e.date)));
            out.println(text("console.valid", text(e.zbs.valid() ? "answer.yes" : "answer.no")));
            if (!Utils.isBlank(e.zbs.notice())) {
                out.println(text("console.note", e.zbs.notice()));
            }
            boolean allow;
            if (e.zbs.invalid()) {
                out.println(text("console.denied"));
                allow = false;
            } else {
                allow = readYesNo(text("console.allow"));
            }
            e.decision = allow;
            if (readYesNo(text("console.persist"))) {
                e.flags = e.flags.with(MF_PERSIST);
            } else {
                e.flags = e.flags.without(MF_PERSIST);
            }
        }
        return pending;
    }

    private boolean readYesNo(String prompt) {
        while (true) {
            out.print(text("console.prompt", prompt));
            out.flush();
            String line;
            try {
                line = in.readLine();
            } catch (Exception e) {
                Logger.error("Console approval read failed: " + e);
                return false;
            }
            if (line == null) {
                return false;
            }
            String s = line.trim().toLowerCase(Locale.ROOT);
            if (s.isEmpty()) {
                continue;
            }
            if (s.startsWith("y")) {
                return true;
            }
            if (s.startsWith("n")) {
                return false;
            }
            out.println(text("console.retry"));
        }
    }

    private static String formatDate(Date date) {
        if (date == null) {
            return text("status.unknown");
        }
        return new SimpleDateFormat(DATE_FORMAT, Locale.ROOT).format(date);
    }
}
