package output

func PrintBanner(p *Printer, toolName, toolVersion string) {
	ascii := `
   ____                 __     __            __      
  / __/__ ___ _________/ /    / /  ___ ___ _/ /__ ___
 _\ \/ -_) _  / __/ __/ _ \  / /__/ -_) _  /  '_/(_-<
/___/\__/\_,_/_/  \__/_//_/ /____/\__/\_,_/_/\_\/___/
                                                     
`
	p.Printf("%s", ascii)
	p.Printf(" haltman.io (https://github.com/haltman-io)\n\n")
	p.Printf(" [codename: %s] - [release: %s]\n\n", toolName, toolVersion)
}
