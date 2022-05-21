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

const routes: Routes = [
  { path: '', component: LandingPageComponent },
  { path: 'register', component: RegistrationPageComponent },
  { path: 'login', component: LoginPageComponent },
  { path: 'companies', component: CompaniesComponent },
  { path: 'user-profile', component: UserProfileComponent },
  { path: 'create-company', component: CreateCompanyComponent },
  { path: 'client-companies', component: ClientCompaniesComponent },
  { path: 'forum', component: ForumComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
