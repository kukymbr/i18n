package i18n_test

type tplData struct {
	TestN int
}

type testStruct struct {
	Test1          string `i18n:"test_1"`
	Test2          string `i18n:"test_2"`
	Test3          string `i18n:"test_3"`
	TestSkip       string `i18n:"-"`
	TestWithoutTag string
}
