import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
	selector: 'app-logout-view',
	templateUrl: './logout-view.component.html',
	styleUrls: ['./logout-view.component.css']
})
export class LogoutViewComponent implements OnInit {

	constructor(private auth: AuthService, private router: Router) { }

	ngOnInit(): void {
		this.auth.logout().subscribe(_ => {
			this.auth.loadStatus()
			this.router.navigate(["/login"])
		})
	}

}
