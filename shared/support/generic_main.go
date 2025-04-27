package shared_support

type GenericMainFunc func()

func GenericMain(mainFunc GenericMainFunc, version, commit, buildDate, builtBy, goLangVersion string) {
	SetVersionBuildInfo(
		version,
		commit,
		buildDate,
		builtBy,
		goLangVersion)

	defer ResetSignals()
	SetupSignalsCatching()

	SetupLogger(StopSignal.Context)

	// Report some information about the build
	Logger().Info().Msgf(
		"%s v%s (git %s, built on %s by %s using go v%s)",
		GetExecutableName(),
		GetVersion(),
		GetCommit(),
		GetBuildDate(),
		GetBuiltBy(),
		GetGoLangVersion())
	if DelveEnabled {
		Logger().Info().Msg("Built with Delve enabled, reporting DEBUG and TRACE logs")
	}

	// Run the main function
	WrapWithSentry(SentryWrappedFunc(mainFunc))
}
