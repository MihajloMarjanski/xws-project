package com.example.agent.service;

import com.example.agent.model.Admin;
import com.example.agent.model.Client;
import com.example.agent.model.CompanyOwner;
import com.example.agent.repository.AdminRepository;
import com.example.agent.repository.ClientRepository;
import com.example.agent.repository.CompanyOwnerRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;


// Ovaj servis je namerno izdvojen kao poseban u ovom primeru.
// U opstem slucaju UserServiceImpl klasa bi mogla da implementira UserDetailService interfejs.
@Service
public class CustomUserDetailsService implements UserDetailsService {

    @Autowired
    private CompanyOwnerRepository companyOwnerRepository;
    @Autowired
    private AdminRepository adminRepository;
    @Autowired
    private ClientRepository clientRepository;


    // Funkcija koja na osnovu username-a iz baze vraca objekat User-a
    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        CompanyOwner companyOwner = companyOwnerRepository.findByUsername(username);
        Admin admin = adminRepository.findByUsername(username);
        Client client = clientRepository.findByUsername(username);

        if (companyOwner==null && admin==null && client==null)
            throw new UsernameNotFoundException(String.format("No user found with username '%s'.", username));
        else if (companyOwner != null)
            return companyOwner;
        else if (admin != null)
            return admin;
        else
            return client;
    }

}
