import { Pipe, PipeTransform } from '@angular/core';
import { PublicUser } from '../services/auth.service';

@Pipe({
	name: 'adminBeforeRegular'
})
export class AdminBeforeRegularPipe implements PipeTransform {

	transform(items: PublicUser[], ...args: unknown[]): PublicUser[] {
		return items.sort((a, b) => {
			if (a.admin == b.admin) {
				return 0
			} else if (a.admin) {
				return -1
			} else {
				return 1
			}
		});
	}

}
