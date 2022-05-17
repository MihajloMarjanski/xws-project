package com.example.agent.repository;

import com.example.agent.model.Company;
import com.example.agent.model.CompanyOwner;
import org.springframework.data.jpa.repository.JpaRepository;

public interface CompanyRepository extends JpaRepository<Company, Integer> {
    Company findByCompanyOwnerId(Integer id);
}
