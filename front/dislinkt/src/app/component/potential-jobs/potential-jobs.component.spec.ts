import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PotentialJobsComponent } from './potential-jobs.component';

describe('PotentialJobsComponent', () => {
  let component: PotentialJobsComponent;
  let fixture: ComponentFixture<PotentialJobsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PotentialJobsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PotentialJobsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
