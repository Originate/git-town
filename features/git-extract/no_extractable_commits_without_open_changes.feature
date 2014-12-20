Feature: git-extract errors if there are not extractable commits

  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'feature' has no extractable commits."
    And I am still on the "feature" branch