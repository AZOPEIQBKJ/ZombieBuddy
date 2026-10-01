import java.io.FileInputStream;
import se.krka.kahlua.j2se.J2SEPlatform;
import se.krka.kahlua.luaj.compiler.LuaCompiler;
import se.krka.kahlua.vm.KahluaThread;

/** Executes the real options Lua on the game's Kahlua, with a small UI fixture. */
public final class LuaSmoke {
    public static void main(String[] args) throws Exception {
        var platform = J2SEPlatform.getInstance();
        var env = platform.newEnvironment();
        var functions = LuaCompiler.class.getDeclaredField("functions");
        functions.setAccessible(true);
        env.rawset("loadstring", ((LuaCompiler[]) functions.get(null))[0]);
        var thread = new KahluaThread(System.out, platform, env);
        var owner = KahluaThread.class.getDeclaredField("debugOwnerThread");
        owner.setAccessible(true);
        owner.set(thread, Thread.currentThread());
        try (var source = new FileInputStream(args[0])) {
            thread.call(LuaCompiler.loadis(source, "Community Lua regression", env), new Object[0]);
        }
    }
}
