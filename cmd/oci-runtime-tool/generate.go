package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	rspec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/opencontainers/runtime-tools/generate"
	"github.com/opencontainers/runtime-tools/generate/seccomp"
	"github.com/urfave/cli"
)

var generateFlags = []cli.Flag{
	cli.StringFlag{Name: "apparmor", Usage: "specifies the the apparmor profile for the container"},
	cli.StringFlag{Name: "arch", Value: runtime.GOARCH, Usage: "architecture the container is created for"},
	cli.StringSliceFlag{Name: "args", Usage: "command to run in the container"},
	cli.StringSliceFlag{Name: "bind", Usage: "bind mount directories src:dest:(rw,ro)"},
	cli.StringSliceFlag{Name: "cap-add", Usage: "add Linux capabilities"},
	cli.StringSliceFlag{Name: "cap-drop", Usage: "drop Linux capabilities"},
	cli.StringFlag{Name: "cgroup", Usage: "cgroup namespace"},
	cli.StringFlag{Name: "cgroups-path", Usage: "specify the path to the cgroups"},
	cli.StringFlag{Name: "cwd", Value: "/", Usage: "current working directory for the process"},
	cli.BoolFlag{Name: "disable-oom-kill", Usage: "disable OOM Killer"},
	cli.StringSliceFlag{Name: "env", Usage: "add environment variable e.g. key=value"},
	cli.IntFlag{Name: "gid", Usage: "gid for the process"},
	cli.StringSliceFlag{Name: "gidmappings", Usage: "add GIDMappings e.g HostID:ContainerID:Size"},
	cli.StringSliceFlag{Name: "groups", Usage: "supplementary groups for the process"},
	cli.StringFlag{Name: "hostname", Usage: "hostname value for the container"},
	cli.StringFlag{Name: "ipc", Usage: "ipc namespace"},
	cli.StringSliceFlag{Name: "label", Usage: "add annotations to the configuration e.g. key=value"},
	cli.Uint64Flag{Name: "linux-cpu-shares", Usage: "the relative share of CPU time available to the tasks in a cgroup"},
	cli.Uint64Flag{Name: "linux-cpu-period", Usage: "the CPU period to be used for hardcapping (in usecs)"},
	cli.Uint64Flag{Name: "linux-cpu-quota", Usage: "the allowed CPU time in a given period (in usecs)"},
	cli.StringFlag{Name: "linux-cpus", Usage: "CPUs to use within the cpuset (default is to use any CPU available)"},
	cli.Uint64Flag{Name: "linux-mem-kernel-limit", Usage: "kernel memory limit (in bytes)"},
	cli.Uint64Flag{Name: "linux-mem-kernel-tcp", Usage: "kernel memory limit for tcp (in bytes)"},
	cli.Uint64Flag{Name: "linux-mem-limit", Usage: "memory limit (in bytes)"},
	cli.Uint64Flag{Name: "linux-mem-reservation", Usage: "memory reservation or soft limit (in bytes)"},
	cli.Uint64Flag{Name: "linux-mem-swap", Usage: "total memory limit (memory + swap) (in bytes)"},
	cli.Uint64Flag{Name: "linux-mem-swappiness", Usage: "how aggressive the kernel will swap memory pages (Range from 0 to 100)"},
	cli.StringFlag{Name: "linux-mems", Usage: "list of memory nodes in the cpuset (default is to use any available memory node)"},
	cli.Int64Flag{Name: "linux-pids-limit", Usage: "maximum number of PIDs"},
	cli.Uint64Flag{Name: "linux-realtime-period", Usage: "CPU period to be used for realtime scheduling (in usecs)"},
	cli.Uint64Flag{Name: "linux-realtime-runtime", Usage: "the time realtime scheduling may use (in usecs)"},
	cli.StringFlag{Name: "mount", Usage: "mount namespace"},
	cli.StringFlag{Name: "mount-cgroups", Value: "no", Usage: "mount cgroups (rw,ro,no)"},
	cli.StringFlag{Name: "mount-label", Usage: "selinux mount context label"},
	cli.StringFlag{Name: "network", Usage: "network namespace"},
	cli.BoolFlag{Name: "no-new-privileges", Usage: "set no new privileges bit for the container process"},
	cli.IntFlag{Name: "oom-score-adj", Usage: "oom_score_adj for the container"},
	cli.StringFlag{Name: "os", Value: runtime.GOOS, Usage: "operating system the container is created for"},
	cli.StringFlag{Name: "output", Usage: "output file (defaults to stdout)"},
	cli.StringFlag{Name: "pid", Usage: "pid namespace"},
	cli.StringSliceFlag{Name: "poststart", Usage: "set command to run in poststart hooks"},
	cli.StringSliceFlag{Name: "poststop", Usage: "set command to run in poststop hooks"},
	cli.StringSliceFlag{Name: "prestart", Usage: "set command to run in prestart hooks"},
	cli.BoolFlag{Name: "privileged", Usage: "enable privileged container settings"},
	cli.BoolFlag{Name: "read-only", Usage: "make the container's rootfs read-only"},
	cli.StringFlag{Name: "root-propagation", Usage: "mount propagation for root"},
	cli.StringFlag{Name: "rootfs", Value: "rootfs", Usage: "path to the rootfs"},
	cli.StringFlag{Name: "seccomp-allow", Usage: "specifies syscalls to respond with allow"},
	cli.StringFlag{Name: "seccomp-arch", Usage: "specifies additional architectures permitted to be used for system calls"},
	cli.StringFlag{Name: "seccomp-default", Usage: "specifies default action to be used for system calls and removes existing rules with specified action"},
	cli.StringFlag{Name: "seccomp-default-force", Usage: "same as seccomp-default but does not remove existing rules with specified action"},
	cli.StringFlag{Name: "seccomp-errno", Usage: "specifies syscalls to respond with errno"},
	cli.StringFlag{Name: "seccomp-kill", Usage: "specifies syscalls to respond with kill"},
	cli.BoolFlag{Name: "seccomp-only", Usage: "specifies to export just a seccomp configuration file"},
	cli.StringFlag{Name: "seccomp-remove", Usage: "specifies syscalls to remove seccomp rules for"},
	cli.BoolFlag{Name: "seccomp-remove-all", Usage: "removes all syscall rules from seccomp configuration"},
	cli.StringFlag{Name: "seccomp-trace", Usage: "specifies syscalls to respond with trace"},
	cli.StringFlag{Name: "seccomp-trap", Usage: "specifies syscalls to respond with trap"},
	cli.StringFlag{Name: "selinux-label", Usage: "process selinux label"},
	cli.StringSliceFlag{Name: "sysctl", Usage: "add sysctl settings e.g net.ipv4.forward=1"},
	cli.StringFlag{Name: "template", Usage: "base template to use for creating the configuration"},
	cli.StringSliceFlag{Name: "tmpfs", Usage: "mount tmpfs"},
	cli.BoolFlag{Name: "tty", Usage: "allocate a new tty for the container process"},
	cli.IntFlag{Name: "uid", Usage: "uid for the process"},
	cli.StringSliceFlag{Name: "uidmappings", Usage: "add UIDMappings e.g HostID:ContainerID:Size"},
	cli.StringFlag{Name: "user", Usage: "user namespace"},
	cli.StringFlag{Name: "uts", Usage: "uts namespace"},
}

var generateCommand = cli.Command{
	Name:   "generate",
	Usage:  "generate an OCI spec file",
	Flags:  generateFlags,
	Before: before,
	Action: func(context *cli.Context) error {
		// Start from the default template.
		specgen := generate.New()

		var template string
		if context.IsSet("template") {
			template = context.String("template")
		}
		if template != "" {
			var err error
			specgen, err = generate.NewFromFile(template)
			if err != nil {
				return err
			}
		}

		err := setupSpec(&specgen, context)
		if err != nil {
			return err
		}

		var exportOpts generate.ExportOptions
		exportOpts.Seccomp = context.Bool("seccomp-only")

		if context.IsSet("output") {
			err = specgen.SaveToFile(context.String("output"), exportOpts)
		} else {
			err = specgen.Save(os.Stdout, exportOpts)
		}
		if err != nil {
			return err
		}
		return nil
	},
}

func setupSpec(g *generate.Generator, context *cli.Context) error {
	if context.GlobalBool("host-specific") {
		g.HostSpecific = true
	}

	spec := g.Spec()

	if len(spec.Version) == 0 {
		g.SetVersion(rspec.Version)
	}

	if context.IsSet("hostname") {
		g.SetHostname(context.String("hostname"))
	}

	g.SetPlatformOS(context.String("os"))
	g.SetPlatformArch(context.String("arch"))

	if context.IsSet("label") {
		annotations := context.StringSlice("label")
		for _, s := range annotations {
			pair := strings.Split(s, "=")
			if len(pair) != 2 {
				return fmt.Errorf("incorrectly specified annotation: %s", s)
			}
			g.AddAnnotation(pair[0], pair[1])
		}
	}

	g.SetRootPath(context.String("rootfs"))

	if context.IsSet("read-only") {
		g.SetRootReadonly(context.Bool("read-only"))
	}

	if context.IsSet("uid") {
		g.SetProcessUID(uint32(context.Int("uid")))
	}

	if context.IsSet("gid") {
		g.SetProcessGID(uint32(context.Int("gid")))
	}

	if context.IsSet("selinux-label") {
		g.SetProcessSelinuxLabel(context.String("selinux-label"))
	}

	g.SetProcessCwd(context.String("cwd"))

	if context.IsSet("apparmor") {
		g.SetProcessApparmorProfile(context.String("apparmor"))
	}

	if context.IsSet("no-new-privileges") {
		g.SetProcessNoNewPrivileges(context.Bool("no-new-privileges"))
	}

	if context.IsSet("tty") {
		g.SetProcessTerminal(context.Bool("tty"))
	}

	if context.IsSet("args") {
		g.SetProcessArgs(context.StringSlice("args"))
	}

	if context.IsSet("env") {
		envs := context.StringSlice("env")
		for _, env := range envs {
			g.AddProcessEnv(env)
		}
	}

	if context.IsSet("groups") {
		groups := context.StringSlice("groups")
		for _, group := range groups {
			groupID, err := strconv.Atoi(group)
			if err != nil {
				return err
			}
			g.AddProcessAdditionalGid(uint32(groupID))
		}
	}

	if context.IsSet("cgroups-path") {
		g.SetLinuxCgroupsPath(context.String("cgroups-path"))
	}

	if context.IsSet("mount-label") {
		g.SetLinuxMountLabel(context.String("mount-label"))
	}

	if context.IsSet("sysctl") {
		sysctls := context.StringSlice("sysctl")
		for _, s := range sysctls {
			pair := strings.Split(s, "=")
			if len(pair) != 2 {
				return fmt.Errorf("incorrectly specified sysctl: %s", s)
			}
			g.AddLinuxSysctl(pair[0], pair[1])
		}
	}

	privileged := false
	if context.IsSet("privileged") {
		privileged = context.Bool("privileged")
	}
	g.SetupPrivileged(privileged)

	if context.IsSet("cap-add") {
		addCaps := context.StringSlice("cap-add")
		for _, cap := range addCaps {
			if err := g.AddProcessCapability(cap); err != nil {
				return err
			}
		}
	}

	if context.IsSet("cap-drop") {
		dropCaps := context.StringSlice("cap-drop")
		for _, cap := range dropCaps {
			if err := g.DropProcessCapability(cap); err != nil {
				return err
			}
		}
	}

	needsNewUser := false

	var uidMaps, gidMaps []string

	if context.IsSet("uidmappings") {
		uidMaps = context.StringSlice("uidmappings")
	}

	if context.IsSet("gidmappings") {
		gidMaps = context.StringSlice("gidmappings")
	}

	if len(uidMaps) > 0 || len(gidMaps) > 0 {
		needsNewUser = true
	}

	setupLinuxNamespaces(context, g, needsNewUser)

	if context.IsSet("tmpfs") {
		tmpfsSlice := context.StringSlice("tmpfs")
		for _, s := range tmpfsSlice {
			dest, options, err := parseTmpfsMount(s)
			if err != nil {
				return err
			}
			g.AddTmpfsMount(dest, options)
		}
	}

	mountCgroupOption := context.String("mount-cgroups")
	if err := g.AddCgroupsMount(mountCgroupOption); err != nil {
		return err
	}

	if context.IsSet("bind") {
		binds := context.StringSlice("bind")
		for _, bind := range binds {
			source, dest, options, err := parseBindMount(bind)
			if err != nil {
				return err
			}
			g.AddBindMount(source, dest, options)
		}
	}

	if context.IsSet("prestart") {
		preStartHooks := context.StringSlice("prestart")
		for _, hook := range preStartHooks {
			path, args := parseHook(hook)
			g.AddPreStartHook(path, args)
		}
	}

	if context.IsSet("poststop") {
		postStopHooks := context.StringSlice("poststop")
		for _, hook := range postStopHooks {
			path, args := parseHook(hook)
			g.AddPostStopHook(path, args)
		}
	}

	if context.IsSet("poststart") {
		postStartHooks := context.StringSlice("poststart")
		for _, hook := range postStartHooks {
			path, args := parseHook(hook)
			g.AddPostStartHook(path, args)
		}
	}

	if context.IsSet("root-propagation") {
		rp := context.String("root-propagation")
		if err := g.SetLinuxRootPropagation(rp); err != nil {
			return err
		}
	}

	for _, uidMap := range uidMaps {
		hid, cid, size, err := parseIDMapping(uidMap)
		if err != nil {
			return err
		}

		g.AddLinuxUIDMapping(hid, cid, size)
	}

	for _, gidMap := range gidMaps {
		hid, cid, size, err := parseIDMapping(gidMap)
		if err != nil {
			return err
		}

		g.AddLinuxGIDMapping(hid, cid, size)
	}

	if context.IsSet("disable-oom-kill") {
		g.SetLinuxResourcesDisableOOMKiller(context.Bool("disable-oom-kill"))
	}

	if context.IsSet("oom-score-adj") {
		g.SetLinuxResourcesOOMScoreAdj(context.Int("oom-score-adj"))
	}

	if context.IsSet("linux-cpu-shares") {
		g.SetLinuxResourcesCPUShares(context.Uint64("linux-cpu-shares"))
	}

	if context.IsSet("linux-cpu-period") {
		g.SetLinuxResourcesCPUPeriod(context.Uint64("linux-cpu-period"))
	}

	if context.IsSet("linux-cpu-quota") {
		g.SetLinuxResourcesCPUQuota(context.Uint64("linux-cpu-quota"))
	}

	if context.IsSet("linux-realtime-runtime") {
		g.SetLinuxResourcesCPURealtimeRuntime(context.Uint64("linux-realtime-runtime"))
	}

	if context.IsSet("linux-realtime-period") {
		g.SetLinuxResourcesCPURealtimePeriod(context.Uint64("linux-realtime-period"))
	}

	if context.IsSet("linux-cpus") {
		g.SetLinuxResourcesCPUCpus(context.String("linux-cpus"))
	}

	if context.IsSet("linux-mems") {
		g.SetLinuxResourcesCPUMems(context.String("linux-mems"))
	}

	if context.IsSet("linux-mem-limit") {
		g.SetLinuxResourcesMemoryLimit(context.Uint64("linux-mem-limit"))
	}

	if context.IsSet("linux-mem-reservation") {
		g.SetLinuxResourcesMemoryReservation(context.Uint64("linux-mem-reservation"))
	}

	if context.IsSet("linux-mem-swap") {
		g.SetLinuxResourcesMemorySwap(context.Uint64("linux-mem-swap"))
	}

	if context.IsSet("linux-mem-kernel-limit") {
		g.SetLinuxResourcesMemoryKernel(context.Uint64("linux-mem-kernel-limit"))
	}

	if context.IsSet("linux-mem-kernel-tcp") {
		g.SetLinuxResourcesMemoryKernelTCP(context.Uint64("linux-mem-kernel-tcp"))
	}

	if context.IsSet("linux-mem-swappiness") {
		g.SetLinuxResourcesMemorySwappiness(context.Uint64("linux-mem-swappiness"))
	}

	err := addSeccomp(context, g)
	return err
}

func setupLinuxNamespaces(context *cli.Context, g *generate.Generator, needsNewUser bool) {
	for _, nsName := range generate.Namespaces {
		if !context.IsSet(nsName) && !(needsNewUser && nsName == "user") {
			continue
		}
		nsPath := context.String(nsName)
		if nsPath == "host" {
			g.RemoveLinuxNamespace(nsName)
			continue
		}
		g.AddOrReplaceLinuxNamespace(nsName, nsPath)
	}
}

func parseIDMapping(idms string) (uint32, uint32, uint32, error) {
	idm := strings.Split(idms, ":")
	if len(idm) != 3 {
		return 0, 0, 0, fmt.Errorf("idmappings error: %s", idms)
	}

	hid, err := strconv.Atoi(idm[0])
	if err != nil {
		return 0, 0, 0, err
	}

	cid, err := strconv.Atoi(idm[1])
	if err != nil {
		return 0, 0, 0, err
	}

	size, err := strconv.Atoi(idm[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return uint32(hid), uint32(cid), uint32(size), nil
}

func parseHook(s string) (string, []string) {
	parts := strings.Split(s, ":")
	args := []string{}
	path := parts[0]
	if len(parts) > 1 {
		args = parts[1:]
	}
	return path, args
}

func parseTmpfsMount(s string) (string, []string, error) {
	var dest string
	var options []string
	var err error

	parts := strings.Split(s, ":")
	if len(parts) == 2 {
		dest = parts[0]
		options = strings.Split(parts[1], ",")
	} else if len(parts) == 1 {
		dest = parts[0]
		options = []string{"rw", "noexec", "nosuid", "nodev", "size=65536k"}
	} else {
		err = fmt.Errorf("invalid value for --tmpfs")
	}

	return dest, options, err
}

func parseBindMount(s string) (string, string, string, error) {
	var source, dest string
	options := "ro"

	bparts := strings.SplitN(s, ":", 3)
	switch len(bparts) {
	case 2:
		source, dest = bparts[0], bparts[1]
	case 3:
		source, dest, options = bparts[0], bparts[1], bparts[2]
	default:
		return source, dest, options, fmt.Errorf("--bind should have format src:dest:[options]")
	}

	return source, dest, options, nil
}

func addSeccomp(context *cli.Context, g *generate.Generator) error {

	// Set the DefaultAction of seccomp
	if context.IsSet("seccomp-default") {
		seccompDefault := context.String("seccomp-default")
		err := g.SetDefaultSeccompAction(seccompDefault)
		if err != nil {
			return err
		}
	} else if context.IsSet("seccomp-default-force") {
		seccompDefaultForced := context.String("seccomp-default-force")
		err := g.SetDefaultSeccompActionForce(seccompDefaultForced)
		if err != nil {
			return err
		}
	}

	// Add the additional architectures permitted to be used for system calls
	if context.IsSet("seccomp-arch") {
		seccompArch := context.String("seccomp-arch")
		architectureArgs := strings.Split(seccompArch, ",")
		for _, arg := range architectureArgs {
			err := g.SetSeccompArchitecture(arg)
			if err != nil {
				return err
			}
		}
	}

	if context.IsSet("seccomp-errno") {
		err := seccompSet(context, "errno", g)
		if err != nil {
			return err
		}
	}

	if context.IsSet("seccomp-kill") {
		err := seccompSet(context, "kill", g)
		if err != nil {
			return err
		}
	}

	if context.IsSet("seccomp-trace") {
		err := seccompSet(context, "trace", g)
		if err != nil {
			return err
		}
	}

	if context.IsSet("seccomp-trap") {
		err := seccompSet(context, "trap", g)
		if err != nil {
			return err
		}
	}

	if context.IsSet("seccomp-allow") {
		err := seccompSet(context, "allow", g)
		if err != nil {
			return err
		}
	}

	if context.IsSet("seccomp-remove") {
		seccompRemove := context.String("seccomp-remove")
		err := g.RemoveSeccompRule(seccompRemove)
		if err != nil {
			return err
		}
	}

	if context.IsSet("seccomp-remove-all") {
		err := g.RemoveAllSeccompRules()
		if err != nil {
			return err
		}
	}
	return nil
}

func seccompSet(context *cli.Context, seccompFlag string, g *generate.Generator) error {
	flagInput := context.String("seccomp-" + seccompFlag)
	flagArgs := strings.Split(flagInput, ",")
	setSyscallArgsSlice := []seccomp.SyscallOpts{}
	for _, flagArg := range flagArgs {
		comparisonArgs := strings.Split(flagArg, ":")
		if len(comparisonArgs) == 5 {
			setSyscallArgs := seccomp.SyscallOpts{
				Action:   seccompFlag,
				Syscall:  comparisonArgs[0],
				Index:    comparisonArgs[1],
				Value:    comparisonArgs[2],
				ValueTwo: comparisonArgs[3],
				Operator: comparisonArgs[4],
			}
			setSyscallArgsSlice = append(setSyscallArgsSlice, setSyscallArgs)
		} else if len(comparisonArgs) == 1 {
			setSyscallArgs := seccomp.SyscallOpts{
				Action:  seccompFlag,
				Syscall: comparisonArgs[0],
			}
			setSyscallArgsSlice = append(setSyscallArgsSlice, setSyscallArgs)
		} else {
			return fmt.Errorf("invalid syscall argument formatting %v", comparisonArgs)
		}

		for _, r := range setSyscallArgsSlice {
			err := g.SetSyscallAction(r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
