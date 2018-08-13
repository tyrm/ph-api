package oauth

import (
	"fmt"
	"net/http"
)

const svgPup = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 90 90" enable-background="new 0 0 90 90" xml:space="preserve"><path d="M84.8,29.3c-2.4-2.3-5.5-10.4-5.5-10.4l-10-4.3h-4.8C60.8,12.8,51.9,8.7,48.2,8c-0.6-0.1-2.4,0-2.4,0s-1.8-0.1-2.4,0  c-3.7,0.7-12.6,4.8-16.3,6.6c0,0,0,0,0,0l-8.6,1.8L6.1,24.7c0,0-3.1,8.4-4.7,9.3c-1.6,0.9-1.2,5.9-1.5,6.4  C-0.4,41,1.4,52.4,1.4,52.4s17.8,9.8,18.5,13.9c0,0-1.3-26.3,0-28.7c0.4,12.6,5.3,23.1,5.4,26.7c0.1,3.6,1.1,15,3.9,16.7  s13.6,1.6,16.6-1.4c3,3,13.9,3.1,16.6,1.4c2.7-1.7,3.7-13.1,3.9-16.7c0.1-3.2,4-11.9,5.1-22.7c0,0,0,0,0,0  c-0.7,16.4-1.5,12.8,0.3,19.2c2.3,7.8,3.8,6.7,3.8,6.7s-0.5-6.1,4.2-9.5c3.1-2.2,6-4.2,6-4.2l4.4-12.5  C90.1,41.2,87.2,31.6,84.8,29.3z M29.4,45c-2.7,0-5.7-3.5-5.7-3.5s3.3-3.5,5.7-3.5c1.9,0,5.1,4.2,5.1,4.2S32.1,45,29.4,45z   M47.4,75.5c-1.6,0-1.6,0-1.6,0s0,0-1.6,0s-4.4-6.7-4-8.5c0.4-1.7,3.5-2.5,5-2.5h0.6h0.6c1.5,0,4.6,0.8,5,2.5  C51.7,68.8,48.9,75.5,47.4,75.5z M62.2,45c-2.7,0-5.1-2.8-5.1-2.8s3.3-4.2,5.1-4.2c2.4,0,5.7,3.5,5.7,3.5S64.9,45,62.2,45z"></path><path d="M45.5,41.9c-2.2,2.2-8.9,2.9-13.3,2.3l0.2,0c0,0,6.2,3.8,12.4,3.8c6.8,0,13-3.7,13-3.7C53.5,44.7,47.6,44,45.5,41.9z"></path></svg>`
const svgHaus = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="-3.456 0 96 124.7725" enable-background="new -3.456 0 96 99.818" xml:space="preserve"><g><path d="M19.91,4.684V93.71h18.547V78.872h14.837V93.71h18.548V4.684H19.91z M34.646,40.381h-7.419V29.252h7.419V40.381z    M49.485,40.381h-7.419V29.252h7.419V40.381z M64.322,40.381h-7.419V29.252h7.419V40.381z"/></g></svg>`


func HandleSVGPup(w http.ResponseWriter, r *http.Request) {
	outputSVG(w, r, svgPup, "Stefania Bonacasa", "https://thenounproject.com/bonste/")
	return
}
func HandleSVGHaus(w http.ResponseWriter, r *http.Request) {
	outputSVG(w, r, svgHaus, "Bjorn Andersson", "https://thenounproject.com/bjorna1/")
	return
}

func outputSVG(w http.ResponseWriter, req *http.Request, svg string, artist string, link string) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("X-Attribution-Artist", artist)
	w.Header().Set("X-Attribution-Link", link)
	w.WriteHeader(200)

	fmt.Fprint(w, svg)
	return
}
