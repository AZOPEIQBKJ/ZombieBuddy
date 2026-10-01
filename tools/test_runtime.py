"""Isolated JVM regression checks. Reads the installed game; never runs its main loop."""
import argparse
import hashlib
import json
import os
import re
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]
PZ_SHA256 = "e1a69eb743ede60b213a0fe7f8b83d4fcab773036d256cc4543a336f3b058a33"
NATIVE_SHA256 = "c2ae9335e717ee24b2f4a40d1a3bf77f1519762a72a0459e766a2bbafc077f6c"


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def literal(text):
    delimiter = "="
    while "]" + delimiter + "]" in text:
        delimiter += "="
    return "[" + delimiter + "[" + text + "]" + delimiter + "]"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--game", type=Path, required=True)
    parser.add_argument("--jdk", type=Path, required=True)
    parser.add_argument("--baseline", type=Path, required=True)
    parser.add_argument("--candidate", type=Path, default=ROOT / "java/build/libs/ZombieBuddy.jar")
    args = parser.parse_args()
    game = args.game.resolve()
    candidate, baseline = args.candidate.resolve(), args.baseline.resolve()
    gamejar = game / "projectzomboid.jar"
    assert digest(gamejar) == PZ_SHA256, "Re-audit the engine before changing its accepted hash"
    before = {p: digest(p) for p in (gamejar, candidate, baseline)}
    reports = ROOT / "artifacts/runtime-tests"
    reports.mkdir(parents=True, exist_ok=True)
    executable = ".exe" if os.name == "nt" else ""
    java = game / ("jre64/bin/java" + executable)
    assert java.is_file(), "Use the audited game's bundled JRE"
    javac = args.jdk / ("bin/javac" + executable)
    results = []
    with tempfile.TemporaryDirectory(prefix="community-runtime-", dir=reports) as tmp:
        work = Path(tmp)
        classes = work / "classes"
        classes.mkdir()
        subprocess.run([str(javac), "-encoding", "UTF-8", "-proc:none", "-cp",
                        os.pathsep.join(map(str, (baseline, gamejar))), "-d", str(classes),
                        str(ROOT / "tools/NativeListSmoke.java"), str(ROOT / "tools/LuaSmoke.java")],
                       check=True, timeout=60)

        def run(name, jar, duplicate=False, native=False):
            case = work / name
            case.mkdir()
            runtimejar = jar
            if native:
                runtimejar = case / "ZombieBuddy.jar"
                shutil.copyfile(jar, runtimejar)
                agent = "-agentpath:" + str(game / "zbNative.dll")
            else:
                agent = "-javaagent:" + str(jar)
            agents = [agent + "=policy=deny-new,verbosity=0,config_dir=" + str(case / "config")]
            if duplicate:
                agents.append(agent + "=policy=allow-all")
            result = subprocess.run([str(java), "-Djava.awt.headless=true", "-Duser.home=" + str(case),
                                     *agents, "-cp", os.pathsep.join(map(str, (classes, runtimejar, gamejar))),
                                     "NativeListSmoke", "baseline" if jar == baseline else "community"],
                                    cwd=case, capture_output=True, text=True,
                                    encoding="utf-8", errors="replace", timeout=90)
            log = result.stdout + result.stderr
            (reports / (name + ".log")).write_text(log, encoding="utf-8")
            assert result.returncode == 0, log[-6000:]
            for error in ("Error transforming", "ERROR applying", "Exception in thread"):
                assert error not in log, log[-6000:]
            calls = log.count("ZomboidFileSystem.loadMods(0 mods) ...")
            expected = 0 if jar == baseline else 3
            assert calls == expected, (name, calls, log[-6000:])
            if duplicate:
                assert "ignoring duplicate" in log
            assert "approval policy unchanged" in log
            results.append(dict(name=name, hookCalls=calls, expected=expected, passed=True))

        run("baseline-misses-List", baseline)
        run("community-integrated-List", candidate)
        run("community-duplicate-keeps-first-policy", candidate, True)
        if os.name == "nt":
            assert digest(game / "zbNative.dll") == NATIVE_SHA256, "Unaudited native DLL"
            before[game / "zbNative.dll"] = NATIVE_SHA256
            run("community-native-integrated-List", candidate, native=True)

        shutil.copyfile(game / "stdlib.lua", work / "stdlib.lua")
        lua = (ROOT / "42/media/lua/client/ZombieBuddy_Options.lua").read_text(encoding="utf-8")
        notification = (ROOT / "42/media/lua/client/ZombieBuddy.lua").read_text(encoding="utf-8")
        translation_root = ROOT / "common/media/lua/shared/Translate"
        english = json.loads((translation_root / "EN/UI.json").read_text(encoding="utf-8"))
        required_keys = set(re.findall(r'getText\("([^"]+)"', lua + notification))
        assert required_keys <= english.keys(), "Missing English UI key"
        assert not any((translation_root / "FR").glob("*")), "French resources must not be shipped"
        fixture = '''
local applied, event = nil, nil
local option = { getValue = function() return 0.4 end }
local options = {
    addSlider = function() return option end,
    addTickBox = function() return option end,
}
PZAPI = { ModOptions = { create = function() return options end } }
Events = { OnMainMenuEnter = { Add = function(fn) event = fn end } }
getText = function(key) return key end
ZombieBuddy = nil
assert(loadstring(SOURCE))()
assert(event == options.apply)
event()
options.apply()
local calls = 0
local count = function(value) assert(value == 0.4); calls = calls + 1 end
ZombieBuddy = {setAutoFixModOrder=count, setFixApprovalDialogCursor=count,
    Watermark={setAlpha=count}, Patches={GameLoadingState={setSuppressSandboxLog=count}}}
event()
options.apply()
assert(calls == 8, "settings did not reach all four existing callbacks twice")
ZombieBuddy = {}
options.apply()
print("PASS options without agent and all existing option callbacks")
local shown = 0
local noop = function() end
getCore = function() return {
    getOptionFontSizeReal=function() return 1 end,
    getScreenWidth=function() return 1280 end,
    getScreenHeight=function() return 720 end,
} end
ISModalRichText = { new=function(self,x,y,w,h,message)
    assert(message == "UI_ZBC_InstallMissing")
    shown = shown + 1
    return {initialise=noop,setY=noop,setVisible=noop,addToUIManager=noop,
        getHeight=function() return h end,chatText={paginate=noop}}
end }
ZombieBuddy = nil
assert(loadstring(NOTIFICATION))()
event()
event()
assert(shown == 1, "missing-agent notification must appear once")
ZombieBuddy = { getVersion=function() return "2.3.3-community.1" end }
assert(loadstring(NOTIFICATION))()
event()
assert(shown == 1, "loaded agent must not show installation notification")
print("PASS missing-agent notification preserved and no upstream installer link")
'''
        test = work / "options.lua"
        test.write_text("local SOURCE=" + literal(lua) + "\nlocal NOTIFICATION=" + literal(notification)
                        + "\n" + fixture, encoding="utf-8")
        result = subprocess.run([str(java), "-Duser.home=" + str(work), "-cp",
                                 os.pathsep.join(map(str, (classes, gamejar))), "LuaSmoke", str(test)],
                                cwd=work, capture_output=True, text=True, encoding="utf-8", timeout=45)
        log = result.stdout + result.stderr
        (reports / "options-lua.log").write_text(log, encoding="utf-8")
        assert result.returncode == 0 and "PASS options" in log, log[-6000:]
        results.append(dict(name="options-Kahlua", passed=True))
    assert all(digest(path) == value for path, value in before.items()), "An input changed during the check"
    (reports / "results.json").write_text(json.dumps(dict(gameSha256=PZ_SHA256,
        candidateSha256=before[candidate], results=results, gameplayAcceptance=False), indent=2) + "\n")
    print(json.dumps(results))


if __name__ == "__main__":
    main()
