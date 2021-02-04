import { Pipe, PipeTransform } from '@angular/core';
import { FmFile } from '../services/filemanager.service';

@Pipe({
	name: 'isdir'
})
export class IsdirPipe implements PipeTransform {

	transform(items: FmFile[], ...args: unknown[]): FmFile[] {
		return items.filter(item => item.isDir);
	}

}
