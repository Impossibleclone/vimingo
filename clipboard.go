package main

type Clipboard struct {
	stack []string

}

func (clip *Clipboard) copy(text string){
	clip.append(clip.stack, text)
}

func (clip *Clipboard) paste() (string){
	if len(clip.stack) == 0{
		return ""
	}else{
		copied := clip.stack[len(clip.stack-1)]
		clip.stack = clip.stack[len(clip.stack-1)]
	}

