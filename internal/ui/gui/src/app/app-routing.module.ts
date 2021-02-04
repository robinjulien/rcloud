import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DashboardViewComponent } from './pages/dashboard/dashboard-view/dashboard-view.component';
import { LoginViewComponent } from './pages/login/login-view/login-view.component';
import { LogoutViewComponent } from './pages/logout/logout-view/logout-view.component';

const routes: Routes = [
	{path: "", component: DashboardViewComponent},
	{path: "login", component: LoginViewComponent },
	{path: "logout", component: LogoutViewComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
