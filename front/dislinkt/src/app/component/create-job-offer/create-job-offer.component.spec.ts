import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateJobOfferComponent } from './create-job-offer.component';

describe('CreateJobOfferComponent', () => {
  let component: CreateJobOfferComponent;
  let fixture: ComponentFixture<CreateJobOfferComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CreateJobOfferComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateJobOfferComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
