package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.3"

func PrintVersion() {
	fmt.Printf("Current subdog version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `                 __         __             
   _____ __  __ / /_   ____/ /____   ____ _
  / ___// / / // __ \ / __  // __ \ / __  /
 (__  )/ /_/ // /_/ // /_/ // /_/ // /_/ / 
/____/ \__,_//_.___/ \__,_/ \____/ \__, /  
                                  /____/
`
	fmt.Printf("%s\n%50s\n\n", banner, "Current subdog version "+version)
}
