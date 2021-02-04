import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
	selector: 'app-dashboard-view',
	templateUrl: './dashboard-view.component.html',
	styleUrls: ['./dashboard-view.component.css']
})
export class DashboardViewComponent implements OnInit {

	constructor(private auth: AuthService, private router: Router) { }

	ngOnInit(): void {
		this.auth.amILoggedIn().subscribe(res => {
			if (res.loggedIn) {
				// Charger fichiers et tout
			} else {
				this.router.navigate(["/login"])
			}
		})
	}

}
