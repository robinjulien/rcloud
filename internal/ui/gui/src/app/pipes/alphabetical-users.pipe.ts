import { Pipe, PipeTransform } from '@angular/core';
import { PublicUser } from '../services/auth.service';

@Pipe({
	name: 'alphabeticalUsers'
})
export class AlphabeticalUsersPipe implements PipeTransform {

	transform(items: PublicUser[], ...args: unknown[]): PublicUser[] {
		return items.sort((a, b) => { return a.id.toUpperCase().localeCompare(b.id.toUpperCase()) });;
	}

}
