import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AdminViewComponent } from './pages/admin/admin-view/admin-view.component';
import { ChangePasswordViewComponent } from './pages/change-password/change-password-view/change-password-view.component';
import { DashboardViewComponent } from './pages/dashboard/dashboard-view/dashboard-view.component';
import { LoginViewComponent } from './pages/login/login-view/login-view.component';
import { LogoutViewComponent } from './pages/logout/logout-view/logout-view.component';

const routes: Routes = [
	{path: "", component: DashboardViewComponent},
	{path: "login", component: LoginViewComponent },
	{path: "logout", component: LogoutViewComponent},
	{path: "admin", component: AdminViewComponent},
	{path: "change-password", component: ChangePasswordViewComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
