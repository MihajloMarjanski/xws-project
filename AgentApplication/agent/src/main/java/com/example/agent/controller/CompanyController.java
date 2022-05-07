package com.example.agent.controller;

import com.example.agent.model.Company;
import com.example.agent.model.CompanyOwner;
import com.example.agent.service.CompanyService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/company")
public class CompanyController {

    @Autowired
    private CompanyService companyService;


    @PostMapping(path = "/createOwner")
    public ResponseEntity<?> createCompanyOwner(@RequestBody CompanyOwner companyOwner) {
        return companyService.saveCompanyOwner(companyOwner);
    }

    @PostMapping(path = "/create")
    public ResponseEntity<?> createCompany(@RequestBody Company company) {
        return companyService.sendCompanyRegistrationRequest(company);
    }
}
