namespace scheduler_service.DataAccess.Entities;

public class Reserving
{
    public long Id { get; set; }

    public string Name { get; set; } = null!;

    public int Age { get; set; }

    public string Email { get; set; } = null!;

    public int PhoneNumber { get; set; }

    public virtual Reservation Reservation { get; set; } = null!;
}