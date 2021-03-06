import { Component, OnInit } from '@angular/core';
import { AuthService } from 'src/app/services/auth.service';

@Component({
	selector: 'app-navbar',
	templateUrl: './navbar.component.html',
	styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {

	constructor(private auth: AuthService) { }

	ngOnInit(): void {
	}

	isLoggedIn(): boolean {
		return this.auth.isLoggedIn()
	}

	getUserID(): string {
		return this.auth.getUser().id
	}

	getUserAdmin(): boolean {
		return this.auth.getUser().admin
	}
}
