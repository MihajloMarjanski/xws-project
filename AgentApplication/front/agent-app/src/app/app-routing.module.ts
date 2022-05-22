import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ClientCompaniesComponent } from './component/client-companies/client-companies.component';
import { CompaniesComponent } from './component/companies/companies.component';
import { CreateCompanyComponent } from './component/create-company/create-company.component';
import { ForumComponent } from './component/forum/forum.component';
import { LandingPageComponent } from './component/landing-page/landing-page.component';
import { LoginPageComponent } from './component/login-page/login-page.component';
import { RegistrationPageComponent } from './component/registration-page/registration-page.component';
import { UserProfileComponent } from './component/user-profile/user-profile.component';
import { AuthGuard } from './service/auth.guard';

const routes: Routes = [
  { path: '', component: LandingPageComponent },
  { path: 'register', component: RegistrationPageComponent },
  { path: 'login', component: LoginPageComponent },
  { path: 'companies', component: CompaniesComponent, canActivate: [AuthGuard], data: { role: ['ROLE_ADMIN']} },
  { path: 'user-profile', component: UserProfileComponent, canActivate: [AuthGuard], data: { role: ['ROLE_CLIENT', 'ROLE_POTENTIAL_OWNER', 'ROLE_OWNER', 'ROLE_ADMIN']} },
  { path: 'create-company', component: CreateCompanyComponent, canActivate: [AuthGuard], data: { role: ['ROLE_POTENTIAL_OWNER']} },
  { path: 'client-companies', component: ClientCompaniesComponent, canActivate: [AuthGuard], data: { role: ['ROLE_CLIENT']} },
  { path: 'forum', component: ForumComponent, canActivate: [AuthGuard], data: { role: ['ROLE_CLIENT']} },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
