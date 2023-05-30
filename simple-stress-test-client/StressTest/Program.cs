using System;
using System.IO;
using System.Net.NetworkInformation;

class Program
{
    static void Main()
    {
        int numPings;
        string numPingsEnv = Environment.GetEnvironmentVariable("num_pings");

        if (!int.TryParse(numPingsEnv, out numPings) || numPings <= 0)
        {
            Console.WriteLine("Invalid value for num_pings environment variable. Using default value: 10");
            numPings = 10;
        }

        string logFile = "/tmp/stress-test-log";
        string host = "www.google.com";
        Ping pingSender = new Ping();

        try
        {
            using (StreamWriter sw = new StreamWriter(logFile))
            {
                Console.SetOut(sw);

                for (int i = 0; i < numPings; i++)
                {
                    PingReply reply = pingSender.Send(host);

                    if (reply.Status == IPStatus.Success)
                    {
                        Console.WriteLine($"Ping {i + 1} to {host} was successful.");
                        Console.WriteLine($"Roundtrip time: {reply.RoundtripTime}ms");
                    }
                    else
                    {
                        Console.WriteLine($"Ping {i + 1} to {host} failed. Status: {reply.Status}");
                    }
                }
            }
        }
        catch (PingException ex)
        {
            using (StreamWriter sw = new StreamWriter(logFile))
            {
                Console.SetOut(sw);
                Console.WriteLine($"An exception occurred while pinging {host}: {ex.Message}");
            }
        }
    }
}