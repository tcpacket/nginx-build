package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/tcpacket/nginx-build/builder"
	"github.com/tcpacket/nginx-build/command"
	"github.com/tcpacket/nginx-build/configure"
	"github.com/tcpacket/nginx-build/logger"
	"github.com/tcpacket/nginx-build/modules"
	"github.com/tcpacket/nginx-build/util"
)

var (
	nginxBuildOptions Options
	log               = logger.Get()
)

func init() {
	nginxBuildOptions = makeNginxBuildOptions()
}

// fake flag for --with-xxx=dynamic
func overrideUnableParseFlags() {
	for i, arg := range os.Args {
		arg = strings.TrimSpace(arg)
		arg = strings.ToLower(arg)
		switch {
		case strings.Contains(arg, "with-http_xslt_module=dynamic"):
			os.Args[i] = "--with-http_xslt_module_dynamic"
		case strings.Contains(arg, "with-http_image_filter_module=dynamic"):
			os.Args[i] = "--with-http_image_filter_module_dynamic"
		case strings.Contains(arg, "with-http_geoip_module=dynamic"):
			os.Args[i] = "--with-http_geoip_module_dynamic"
		case strings.Contains(arg, "with-http_perl_module=dynamic"):
			os.Args[i] = "--with-http_perl_module_dynamic"
		case strings.Contains(arg, "with-mail=dynamic"):
			os.Args[i] = "--with-mail_dynamic"
		case strings.Contains(arg, "with-stream=dynamic"):
			os.Args[i] = "--with-stream_dynamic"
		case strings.Contains(arg, "with-stream_geoip_module=dynamic"):
			os.Args[i] = "--with-stream_geoip_module_dynamic"
		case strings.Contains(arg, "with-compat"):
			os.Args[i] = "--with-compat"
		default:
			continue
		}
	}
}

func main() {
	var (
		multiflagPatch StringFlag
	)

	// Parse flags
	for k, v := range nginxBuildOptions.Bools {
		v.Enabled = flag.Bool(k, false, v.Desc)
		nginxBuildOptions.Bools[k] = v
	}
	for k, v := range nginxBuildOptions.Values {
		if k == "patch" {
			flag.Var(&multiflagPatch, k, v.Desc)
		} else {
			v.Value = flag.String(k, v.Default, v.Desc)
			nginxBuildOptions.Values[k] = v
		}
	}
	for k, v := range nginxBuildOptions.Numbers {
		v.Value = flag.Int(k, v.Default, v.Desc)
		nginxBuildOptions.Numbers[k] = v
	}

	overrideUnableParseFlags()

	var (
		configureOptions configure.Options
		multiflag        StringFlag
		multiflagDynamic StringFlag
	)

	argsString := configure.MakeArgsString()
	for k, v := range argsString {
		if k == "add-module" {
			flag.Var(&multiflag, k, v.Desc)
		} else if k == "add-dynamic-module" {
			flag.Var(&multiflagDynamic, k, v.Desc)
		} else {
			v.Value = flag.String(k, "", v.Desc)
			argsString[k] = v
		}
	}

	argsBool := configure.MakeArgsBool()
	for k, v := range argsBool {
		v.Enabled = flag.Bool(k, false, v.Desc)
		argsBool[k] = v
	}

	flag.CommandLine.SetOutput(os.Stdout)
	// The output of original flag.Usage() is too long
	defaultUsage := flag.Usage
	flag.Usage = usage
	flag.Parse()

	jobs := nginxBuildOptions.Numbers["j"].Value

	verbose := nginxBuildOptions.Bools["verbose"].Enabled
	pcreStatic := nginxBuildOptions.Bools["pcre"].Enabled
	openSSLStatic := nginxBuildOptions.Bools["openssl"].Enabled
	libreSSLStatic := nginxBuildOptions.Bools["libressl"].Enabled
	zlibStatic := nginxBuildOptions.Bools["zlib"].Enabled
	clear := nginxBuildOptions.Bools["clear"].Enabled
	versionPrint := nginxBuildOptions.Bools["version"].Enabled
	versionsPrint := nginxBuildOptions.Bools["versions"].Enabled
	openResty := nginxBuildOptions.Bools["openresty"].Enabled
	tengine := nginxBuildOptions.Bools["tengine"].Enabled
	configureOnly := nginxBuildOptions.Bools["configureonly"].Enabled
	idempotent := nginxBuildOptions.Bools["idempotent"].Enabled
	helpAll := nginxBuildOptions.Bools["help-all"].Enabled
	enableHTTP2 := nginxBuildOptions.Bools["with-http_v2_module"].Enabled
	version := nginxBuildOptions.Values["v"].Value
	nginxConfigurePath := nginxBuildOptions.Values["c"].Value
	modulesConfPath := nginxBuildOptions.Values["m"].Value
	workParentDir := nginxBuildOptions.Values["d"].Value
	pcreVersion := nginxBuildOptions.Values["pcreversion"].Value
	openSSLVersion := nginxBuildOptions.Values["opensslversion"].Value
	libreSSLVersion := nginxBuildOptions.Values["libresslversion"].Value
	zlibVersion := nginxBuildOptions.Values["zlibversion"].Value
	openRestyVersion := nginxBuildOptions.Values["openrestyversion"].Value
	tengineVersion := nginxBuildOptions.Values["tengineversion"].Value
	patchOption := nginxBuildOptions.Values["patch-opt"].Value

	// Allow multiple flags for `--patch`
	{
		tmp := nginxBuildOptions.Values["patch"]
		tmp2 := multiflagPatch.String()
		tmp.Value = &tmp2
		nginxBuildOptions.Values["patch"] = tmp
	}

	// Allow multiple flags for `--add-module`
	{
		tmp := argsString["add-module"]
		tmp2 := multiflag.String()
		tmp.Value = &tmp2
		argsString["add-module"] = tmp
	}

	// Allow multiple flags for `--add-dynamic-module`
	{
		tmp := argsString["add-dynamic-module"]
		tmp2 := multiflagDynamic.String()
		tmp.Value = &tmp2
		argsString["add-dynamic-module"] = tmp
	}

	patchPath := nginxBuildOptions.Values["patch"].Value
	configureOptions.Values = argsString
	configureOptions.Bools = argsBool

	if *helpAll {
		defaultUsage()
		return
	}

	if *versionPrint {
		printNginxBuildVersion()
		return
	}

	if *versionsPrint {
		printNginxVersions()
		return
	}

	printFirstMsg()

	// set verbose mode
	command.VerboseEnabled = *verbose

	var nginxBuilder builder.Builder
	if *openResty && *tengine {
		log.Fatal().Msg("select one between '-openresty' and '-tengine'.")
	}
	if *openSSLStatic && *libreSSLStatic {
		log.Fatal().Msg("select one between '-openssl' and '-libressl'.")
	}
	if *openResty {
		nginxBuilder = builder.MakeBuilder(builder.ComponentOpenResty, *openRestyVersion)
	} else if *tengine {
		nginxBuilder = builder.MakeBuilder(builder.ComponentTengine, *tengineVersion)
	} else {
		nginxBuilder = builder.MakeBuilder(builder.ComponentNginx, *version)
	}
	pcreBuilder := builder.MakeLibraryBuilder(builder.ComponentPcre, *pcreVersion, *pcreStatic)
	openSSLBuilder := builder.MakeLibraryBuilder(builder.ComponentOpenSSL, *openSSLVersion, *openSSLStatic)
	libreSSLBuilder := builder.MakeLibraryBuilder(builder.ComponentLibreSSL, *libreSSLVersion, *libreSSLStatic)
	zlibBuilder := builder.MakeLibraryBuilder(builder.ComponentZlib, *zlibVersion, *zlibStatic)

	if *idempotent {
		builders := []builder.Builder{
			nginxBuilder,
			pcreBuilder,
			openSSLBuilder,
			libreSSLBuilder,
			zlibBuilder,
		}

		isSame, err := builder.IsSameVersion(builders)
		if err != nil {
			log.Warn().Msgf("failed to check version: %s", err)
		}
		if isSame {
			log.Warn().Msgf("Installed nginx is same version.")
			return
		}
	}

	// change default umask
	_ = syscall.Umask(0)

	versionCheck(*version)

	nginxConfigure, err := util.FileGetContents(*nginxConfigurePath)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	nginxConfigure = configure.Normalize(nginxConfigure)

	modules3rd, err := modules.Load(*modulesConfPath)
	if err != nil {
		log.Fatal().Msgf("Failed to load modules: %v", err)
	}

	if len(*workParentDir) == 0 {
		log.Fatal().Msg("set working directory with -d")
	}

	if !util.FileExists(*workParentDir) {
		err := os.Mkdir(*workParentDir, 0755)
		if err != nil {
			log.Fatal().Msgf("Failed to create working directory '%s': %v", *workParentDir, err)
		}
	}

	var workDir string
	if *openResty {
		workDir = *workParentDir + "/openresty/" + *openRestyVersion
	} else if *tengine {
		workDir = *workParentDir + "/tengine/" + *tengineVersion
	} else {
		workDir = *workParentDir + "/nginx/" + *version
	}

	if *clear {
		err := util.ClearWorkDir(workDir)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	}

	if !util.FileExists(workDir) {
		err := os.MkdirAll(workDir, 0755)
		if err != nil {
			log.Fatal().Msgf("Failed to create working directory(%s) does not exist.", workDir)
		}
	}

	rootDir := util.SaveCurrentDir()
	err = os.Chdir(workDir)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}

	// remove nginx source code applyed patch
	if *patchPath != "" && util.FileExists(nginxBuilder.SourcePath()) {
		err := os.RemoveAll(nginxBuilder.SourcePath())
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	}

	var wg sync.WaitGroup
	if *pcreStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&pcreBuilder)
			wg.Done()
		}()
	}

	if *openSSLStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&openSSLBuilder)
			wg.Done()
		}()
	}

	if *libreSSLStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&libreSSLBuilder)
			wg.Done()
		}()
	}

	if *zlibStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&zlibBuilder)
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		downloadAndExtractParallel(&nginxBuilder)
		wg.Done()
	}()

	if len(modules3rd) > 0 {
		wg.Add(len(modules3rd))
		for _, m := range modules3rd {
			go func(m modules.Module) {
				modules.DownloadAndExtractParallel(m)
				wg.Done()
			}(m)
		}

	}

	// wait until all downloading processes by goroutine finish
	wg.Wait()

	if len(modules3rd) > 0 {
		for _, m := range modules3rd {
			if err := modules.Provide(&m); err != nil {
				log.Fatal().Msgf("%v", err)
			}
		}
	}

	// cd workDir/nginx-${version}
	err = os.Chdir(nginxBuilder.SourcePath())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to change directory")
	}

	var dependencies []builder.StaticLibrary
	if *pcreStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&pcreBuilder))
	}

	if *openSSLStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&openSSLBuilder))
	}

	if *libreSSLStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&libreSSLBuilder))
	}

	if *zlibStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&zlibBuilder))
	}

	log.Printf("Generating configure script for %s.....", nginxBuilder.SourcePath())

	if *pcreStatic && pcreBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Warn().Msgf("%v", pcreBuilder.WarnMsgWithLibrary())
	}

	if *openSSLStatic && openSSLBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Warn().Msgf("%v", openSSLBuilder.WarnMsgWithLibrary())
	}

	if *libreSSLStatic && libreSSLBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Warn().Msgf("%v", libreSSLBuilder.WarnMsgWithLibrary())
	}

	if *zlibStatic && zlibBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Warn().Msgf("%v", zlibBuilder.WarnMsgWithLibrary())
	}

	conf := &configure.Conf{
		Mods: modules3rd,
		Deps: dependencies,
		Opts: configureOptions,
		Dir:  rootDir,

		HTTP2Enabled:     *enableHTTP2,
		OpenRestyEnabled: *openResty,
	}

	configureScript := conf.Generate(nginxConfigure, *jobs)

	err = os.WriteFile("./nginx-configure", []byte(configureScript), 0655)
	if err != nil {
		log.Fatal().Msgf("Failed to generate configure script for %s", nginxBuilder.SourcePath())
	}

	util.Patch(*patchPath, *patchOption, rootDir, false)

	// reverts source code with patch -R when the build was interrupted.
	if *patchPath != "" {
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, os.Interrupt)
		go func() {
			<-sigChannel
			util.Patch(*patchPath, *patchOption, rootDir, true)
		}()
	}

	log.Printf("Configure %s.....", nginxBuilder.SourcePath())

	err = configure.Run()
	if err != nil {
		log.Printf("Failed to configure %s\n", nginxBuilder.SourcePath())
		util.Patch(*patchPath, *patchOption, rootDir, true)
		util.PrintFatalMsg(err, "nginx-configure.log")
	}

	if *configureOnly {
		util.Patch(*patchPath, *patchOption, rootDir, true)
		printLastMsg(workDir, nginxBuilder.SourcePath(), *openResty, *configureOnly)
		return
	}

	log.Printf("Build %s.....", nginxBuilder.SourcePath())

	if *openSSLStatic {
		// Sometimes machine hardware name('uname -m') is different
		// from machine processor architecture name('uname -p') on Mac.
		// Specifically, `uname -p` is 'i386' and `uname -m` is 'x86_64'.
		// In this case, a build of OpenSSL fails.
		// So it needs to convince OpenSSL with KERNEL_BITS.

		//goland:noinspection GoBoolExpressions
		if runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" {
			err := os.Setenv("KERNEL_BITS", "64")
			if err != nil {
				log.Fatal().Err(err).Msgf("failed to set kernel_bits environment variable")
			}
		}
	}

	err = builder.BuildNginx(*jobs)
	if err != nil {
		log.Printf("Failed to build %s\n", nginxBuilder.SourcePath())
		util.Patch(*patchPath, *patchOption, rootDir, true)
		util.PrintFatalMsg(err, "nginx-build.log")
	}

	printLastMsg(workDir, nginxBuilder.SourcePath(), *openResty, *configureOnly)
}
