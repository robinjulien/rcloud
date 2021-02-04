import { Pipe, PipeTransform } from '@angular/core';

const videoFilesExt = ["webm", "mp4", "avi", "wmv", "mov", "mkv"]
const imageFilesExt = ["jpg", "png", "gif", "webp", "tiff", "psd", "bmp", "jpeg", "svg"]
const docFilesExt = ["doc", "docx", "xls", "xlsx", "ppt", "pptx", "ppsm", "odt", "ods", "odp", "odg", "odc", "odf", "odb", "odi", "odm", "pdf"]
const musicFilesExt = ["wav", "ogg", "flac", "mp3"]
const codeFilesExt = ["c", "h", "cpp", "hpp", "js", "css", "html", "php", "go", "java", "py", "rs", "ruby", "jl", "f", "for", "cs", "swift", "sql", "sh", "bat", "kt", "kts", "lua", "r"]

@Pipe({
	name: 'exttoimg'
})
export class ExttoimgPipe implements PipeTransform {

	transform(ext: string, ...args: unknown[]): string {
		let icon = "file.svg"

		if (videoFilesExt.includes(ext)) {
			icon = "film.svg"
		} else if (imageFilesExt.includes(ext)) {
			icon =  "image.svg"
		} else if (docFilesExt.includes(ext)) {
			icon = "file-text.svg"
		} else if (musicFilesExt.includes(ext)) {
			icon = "music.svg"
		} else if (codeFilesExt.includes(ext)) {
			icon = "code.svg"
		}

		return icon
	}

}
