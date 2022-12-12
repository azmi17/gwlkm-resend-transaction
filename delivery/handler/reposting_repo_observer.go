package handler

import (
	"fmt"
	"gwlkm-resend-transaction/entities"
	"time"

	"github.com/kpango/glg"
	"github.com/schollz/progressbar/v3"
)

// Temporary Functions
func RepostingSchedulerRepoObserver() {

	var bar *progressbar.ProgressBar
	var numOfSuccess, numOfFailed int

	for po := range entities.PrintRepoChan {
		if po.Status == entities.PRINT_INIT_REPO_CHAN {
			numOfSuccess = 0
			numOfFailed = 0
			bar = progressbar.NewOptions(po.Size,
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionShowBytes(false),
				progressbar.OptionSetWidth(20),
				progressbar.OptionSetDescription("[reset]Reposting saldo apex is processing..."),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}))
		} else if po.Status == entities.PRINT_FINISH_REPO_CHAN {
			fmt.Println("")
			_ = glg.Log("Scheduler INFO: Reposting Success =>", numOfSuccess)
			_ = glg.Log("Scheduler INFO: Reposting Failed =>", numOfFailed)
			_ = glg.Log("Scheduler INFO: Reposting total =>", numOfSuccess+numOfFailed)

			hours, minutes, _ := time.Now().Clock()
			currUTCTimeInString := fmt.Sprintf("%d:%02d", hours, minutes)
			_ = glg.Log("Scheduler INFO:", "Reposting saldo apex is done at:", currUTCTimeInString)
			fmt.Println("")

			bar.Finish()
		} else {
			if po.Status == entities.PRINT_SUCCESS_STATUS_REPO_CHAN {
				numOfSuccess++
				// glg.Log("Reposting Success:", "=> Kode LKM:", po.KodeLKM, ", => Messages:", po.Message)
			} else {
				numOfFailed++
				glg.Fail("Reposting Errors:", "=> Kode LKM:", po.KodeLKM, ", => Messages:", po.Message)
			}
			bar.Add(1)
			// time.Sleep(1 * time.Nanosecond) // debug mode..
		}

	}

}
