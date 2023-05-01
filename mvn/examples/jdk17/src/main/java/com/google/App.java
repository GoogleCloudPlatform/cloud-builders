package com.google;

/**
 * Hello world using the Record feature of jdk17 to ensure that the builder we are using supports at
 * least jdk17.
 *
 */

public class App
{
    record Person(String firstName, String lastName) { }

    public static void main( String[] args )
    {
      Person person = new Person("John", "Doe");
      System.out.println( "Hello " + person.firstName + " " + person.lastName + "!" );
    }
}
