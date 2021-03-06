import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ChatComponent } from './component/chat/chat.component';
import { CreateJobOfferComponent } from './component/create-job-offer/create-job-offer.component';
import { CreatePostComponent } from './component/create-post/create-post.component';
import { LandingPageComponent } from './component/landing-page/landing-page.component';
import { LoginPageComponent } from './component/login-page/login-page.component';
import { NotificationsComponent } from './component/notifications/notifications.component';
import { OffersComponent } from './component/offers/offers.component';
import { PasswordlessComponent } from './component/passwordless/passwordless.component';
import { PostsComponent } from './component/posts/posts.component';
import { PotentialConnectionsComponent } from './component/potential-connections/potential-connections.component';
import { PotentialJobsComponent } from './component/potential-jobs/potential-jobs.component';
import { RegistrationPageComponent } from './component/registration-page/registration-page.component';
import { UserProfileComponent } from './component/user-profile/user-profile.component';
import { UsersComponent } from './component/users/users.component';
import { ViewPostComponent } from './component/view-post/view-post.component';
import { AuthGuard } from './service/auth.guard';

const routes: Routes = [
  { path: '', component: LandingPageComponent },
  { path: 'register', component: RegistrationPageComponent },
  { path: 'login', component: LoginPageComponent },
  { path: 'user-profile', component: UserProfileComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER', 'ROLE_ADMIN']} },
  { path: 'create-job-offer', component: CreateJobOfferComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER']} },
  { path: 'users', component: UsersComponent},
  { path: 'offers', component: OffersComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER']} },
  { path: 'chat', component: ChatComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER']} },
  { path: 'posts', component: PostsComponent},
  { path: 'view-post', component: ViewPostComponent},
  { path: 'notifications', component: NotificationsComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER']} },
  { path: 'create-post', component: CreatePostComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER']} },
  { path: 'passwordless', component: PasswordlessComponent},
  { path: 'recommended-connections', component: PotentialConnectionsComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER']}},
  { path: 'recommended-jobs', component: PotentialJobsComponent, canActivate: [AuthGuard], data: { role: ['ROLE_USER']}},

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
